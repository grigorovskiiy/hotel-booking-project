package handlers

import (
	"encoding/json"
	"github.com/IvanChumakov/hotel-booking-project/internal/app/services"
	"github.com/IvanChumakov/hotel-booking-project/internal/models"
	"io"
	"log"
	"net/http"
)

func GetBookings(w http.ResponseWriter, r *http.Request) {
	log.Print("/get_bookings")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	bookings, err := services.GetAllBookings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(bookings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func GetBookingsByName(w http.ResponseWriter, r *http.Request) {
	log.Print("/get_bookings_by_name")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var hotel models.Hotels
	err = json.Unmarshal(data, &hotel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	bookings, err := services.GetaBookingByName(hotel.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(bookings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func GetFreeRoomsByDate(w http.ResponseWriter, r *http.Request) {
	log.Print("/get_free_rooms_by_date")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var booking models.Booking
	err = json.Unmarshal(data, &booking)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rooms, err := services.GetHotelRoomsWithPrice(booking)
	freeRooms, err := services.FilterRooms(booking, rooms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(freeRooms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AddBooking(w http.ResponseWriter, r *http.Request) {
	log.Print("/add_booking")
	if http.MethodPost != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var booking models.Booking
	err = json.Unmarshal(data, &booking)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = services.MakePaymentOperation(booking)
	if err != nil {
		http.Error(w, "Failed to make payment operation: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = services.SendNotification(booking)
	w.WriteHeader(http.StatusOK)
}

func PaymentCallBack(w http.ResponseWriter, r *http.Request) {
	log.Print("/payment_callback")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var paymentInfo models.PaymentInfo
	err = json.Unmarshal(body, &paymentInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = services.AddBooking(paymentInfo.Booking)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
