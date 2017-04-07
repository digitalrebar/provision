package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/client/reservations"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

type ReservationOps struct{}

func (be ReservationOps) GetType() interface{} {
	return &models.Reservation{}
}

func (be ReservationOps) GetId(obj interface{}) (string, error) {
	reservation, ok := obj.(*models.Reservation)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to reservation create")
	}
	return reservation.Addr.String(), nil
}

func (be ReservationOps) List() (interface{}, error) {
	d, e := session.Reservations.ListReservations(reservations.NewListReservationsParams(), basicAuth)
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
	d, e := session.Reservations.GetReservation(reservations.NewGetReservationParams().WithAddress(s), basicAuth)
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
	d, e := session.Reservations.CreateReservation(reservations.NewCreateReservationParams().WithBody(reservation), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ReservationOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to reservation patch")
	}
	s, e := convertStringToAddress(id)
	if e != nil {
		return nil, e
	}
	d, e := session.Reservations.PatchReservation(reservations.NewPatchReservationParams().WithAddress(s).WithBody(data), basicAuth)
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
	d, e := session.Reservations.DeleteReservation(reservations.NewDeleteReservationParams().WithAddress(s), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addReservationCommands()
	App.AddCommand(tree)
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
