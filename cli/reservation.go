package main

import (
	"fmt"

	"github.com/rackn/rocket-skates/client/reservations"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

type ReservationOps struct{}

func (be ReservationOps) GetType() interface{} {
	return &models.Reservation{}
}

func (be ReservationOps) List() (interface{}, error) {
	d, e := session.Reservations.ListReservations(reservations.NewListReservationsParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ReservationOps) Get(id string) (interface{}, error) {
	s, e := convertStringToAddress(id)
	if e != nil {
		return nil, e
	}
	d, e := session.Reservations.GetReservation(reservations.NewGetReservationParams().WithAddress(s))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ReservationOps) Create(obj interface{}) (interface{}, error) {
	reservation, ok := obj.(*models.Reservation)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to reservation create")
	}
	d, e := session.Reservations.CreateReservation(reservations.NewCreateReservationParams().WithBody(reservation))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ReservationOps) Put(id string, obj interface{}) (interface{}, error) {
	reservation, ok := obj.(*models.Reservation)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to reservation put")
	}
	s, e := convertStringToAddress(id)
	if e != nil {
		return nil, e
	}
	d, e := session.Reservations.PutReservation(reservations.NewPutReservationParams().WithAddress(s).WithBody(reservation))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ReservationOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.([]*models.JSONPatchOperation)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to reservation patch")
	}
	s, e := convertStringToAddress(id)
	if e != nil {
		return nil, e
	}
	d, e := session.Reservations.PatchReservation(reservations.NewPatchReservationParams().WithAddress(s).WithBody(data))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ReservationOps) Delete(id string) (interface{}, error) {
	s, e := convertStringToAddress(id)
	if e != nil {
		return nil, e
	}
	d, e := session.Reservations.DeleteReservation(reservations.NewDeleteReservationParams().WithAddress(s))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addReservationCommands()
	app.AddCommand(tree)
}

func addReservationCommands() (res *cobra.Command) {
	singularName := "reservation"
	name := "reservations"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	commands := commonOps(singularName, name, &ReservationOps{})
	res.AddCommand(commands...)
	return res
}
