package db

import (
	"errors"
	"gorm.io/gorm"
)

type TxWrapper[T any] struct {
	Tx *gorm.DB
}

func NewTxWrapper[T any](tx *gorm.DB) *TxWrapper[T] {
	return &TxWrapper[T]{tx}
}

func (t *TxWrapper[T]) FindOne(conditions ...any) (*T, error) {
	var obj T
	err := t.Tx.Take(&obj, conditions...).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
func (t *TxWrapper[T]) MustFindOne(conditions ...any) *T {
	obj, err := t.FindOne(conditions...)
	if err != nil {
		panic(err)
	}
	return obj
}
func (t *TxWrapper[T]) FindMany(conditions ...any) ([]T, error) {
	var arr []T
	err := t.Tx.Find(&arr, conditions...).Error
	if err != nil {
		return nil, err
	}
	return arr, nil
}
func (t *TxWrapper[T]) MustFindMany(conditions ...any) []T {
	arr, err := t.FindMany(conditions...)
	if err != nil {
		panic(err)
	}
	return arr
}
func (t *TxWrapper[T]) MustCreate(obj *T) {
	err := t.Tx.Create(&obj).Error
	if err != nil {
		panic(err)
	}
}
func (t *TxWrapper[T]) MustCreateMany(arr []T) {
	err := t.Tx.Create(&arr).Error
	if err != nil {
		panic(err)
	}
}
func (t *TxWrapper[T]) MustSave(obj *T) {
	err := t.Tx.Save(&obj).Error
	if err != nil {
		panic(err)
	}
}
func (t *TxWrapper[T]) MustDelete() {
	var empty T
	err := t.Tx.Delete(&empty).Error
	if err != nil {
		panic(err)
	}
}
func (t *TxWrapper[T]) MustUpdates(values any) {
	var err error
	if t.Tx.Statement.Model == nil {
		var empty T
		err = t.Tx.Model(&empty).Updates(values).Error
	} else {
		err = t.Tx.Updates(values).Error
	}
	if err != nil {
		panic(err)
	}
}
func (t *TxWrapper[T]) MustUpdate(column string, value any) {
	var err error
	if t.Tx.Statement.Model == nil {
		var empty T
		err = t.Tx.Model(&empty).Update(column, value).Error
	} else {
		err = t.Tx.Update(column, value).Error
	}
	if err != nil {
		panic(err)
	}
}
func (t *TxWrapper[T]) Where(query any, args ...any) *TxWrapper[T] {
	return NewTxWrapper[T](t.Tx.Where(query, args...))
}
func (t *TxWrapper[T]) Select(query any, args ...any) *TxWrapper[T] {
	return NewTxWrapper[T](t.Tx.Select(query, args...))
}
func (t *TxWrapper[T]) Raw(sql string, values ...any) *TxWrapper[T] {
	return NewTxWrapper[T](t.Tx.Raw(sql, values...))
}
func (t *TxWrapper[T]) MustExec(sql string, values ...any) {
	err := t.Tx.Exec(sql, values...).Error
	if err != nil {
		panic(err)
	}
}
func (t *TxWrapper[T]) Model(value any) *TxWrapper[T] {
	return NewTxWrapper[T](t.Tx.Model(value))
}
func (t *TxWrapper[T]) Order(value any) *TxWrapper[T] {
	return NewTxWrapper[T](t.Tx.Order(value))
}
func (t *TxWrapper[T]) MustScan(dest any) {
	var err error
	if t.Tx.Statement.Model == nil {
		var empty T
		err = t.Tx.Model(&empty).Scan(&dest).Error
	} else {
		err = t.Tx.Scan(&dest).Error
	}
	if err != nil {
		panic(err)
	}
}
