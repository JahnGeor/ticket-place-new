package services

import (
	"fptr/internal/entities"
	apperr "fptr/internal/error_list"
	"fptr/internal/gateways"
	"github.com/google/logger"
	"net/http"
)

type ClientService struct {
	gw *gateways.Gateway
	*logger.Logger
}

func NewClientService(gw *gateways.Gateway, logg *logger.Logger) *ClientService {
	return &ClientService{
		gw:     gw,
		Logger: logg,
	}
}

func (s *ClientService) GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, error) {
	click, err := s.gw.GetLastReceipt(connectionURL, session)
	if err != nil {
		switch err.(type) {
		case *apperr.ClientError:
			if !(err.(*apperr.ClientError).StatusCode == http.StatusNotFound) {
				s.Errorf("Ошибка при запросе последнего заказа: %v", err)
			}
		}

		return nil, err
	}

	return click, nil
}

func (s *ClientService) PrintSell(info entities.Info, id string) error {
	sell, err := s.gw.Listener.GetSell(info, id)
	if err != nil {
		s.Errorf("Ошибка во время печати заказа с номером %s, клиент: %v", id, err)
		return err
	}

	if err = s.gw.KKT.PrintSell(*sell); err != nil {
		switch err.(type) {
		case *apperr.BusinessError:
			s.Warningf("Ошибка во время печати чека продажи заказа с номером %s, ККТ: %v", id, err)
		default:
			s.Errorf("Ошибка во время печати чека продажи заказа с номером %s, ККТ: %v", id, err)
		}
		return err
	}
	s.Infof("Выполнена печать чека заказа с номером: %s\n", id)
	return nil
}

func (s *ClientService) PrintRefoundFromSell(info entities.Info, id string) error {
	sell, err := s.gw.Listener.GetSell(info, id)
	if err != nil {
		s.Errorf("Ошибка во время печати возврата заказа с номером %s, клиент: %v", id, err)
		return err
	}
	err = s.gw.KKT.PrintRefoundFromCheck(*sell)
	if err != nil {
		switch err.(type) {
		case *apperr.BusinessError:
			s.Warningf("Ошибка во время печати возврата заказа с номером %s, ККТ: %v", id, err)
		default:
			s.Errorf("Ошибка во время печати возврата заказа с номером %s, ККТ: %v", id, err)
		}
		return err
	}
	s.Infof("Выполнена печать чека возврата заказа с номером: %s\n", id)
	return nil
}

func (s *ClientService) PrintRefound(info entities.Info, id string) error {
	refound, err := s.gw.Listener.GetRefound(info, id)
	if err != nil {
		s.Errorf("Ошибка во время печати возврата заказа с номером %s, клиент: %v", id, err)
		return err
	}

	err = s.gw.KKT.PrintRefound(*refound)
	if err != nil {
		switch err.(type) {
		case *apperr.BusinessError:
			s.Warningf("Ошибка во время печати возврата заказа с номером %s, ККТ: %v", id, err)
		default:
			s.Errorf("Ошибка во время печати возврата заказа с номером %s, ККТ: %v", id, err)
		}
		return err
	}

	s.Infof("Выполнена печать чека возврата заказа с номером: %s\n", id)
	return nil
}

func (s *ClientService) Login(config entities.AppConfig) (*entities.SessionInfo, error) {
	session, err := s.gw.Login(config)
	if err != nil {
		s.Errorf("Во время авторизации произошла ошибка: %v", err)
		return nil, err
	}
	return session, nil
}
