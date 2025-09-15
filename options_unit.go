package kboot

type UnitOption Option[*unitImpl]

func DependsOn(dep ...string) UnitOption {
	return optionFunc[*unitImpl](func(unit *unitImpl) {
		unit.depends = append(unit.depends, dep...)
	})
}
