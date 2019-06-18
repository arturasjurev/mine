package factory

type Factory struct {
	// Factory equipment.
	Grinder Equipment
	Freezer Equipment
	Smelter Equipment

	// Resource factory job.
	Resource Mineral
	From     MineralState
	To       MineralState
}

func (f *Factory) Grind() error {
	if err := f.Grinder.Insert(f.Resource); err != nil {
		return err
	}
	if err := f.Grinder.Process(); err != nil {
		return err
	}
	if product, err := f.Grinder.Takeout(); err != nil {
		return err
	} else {
		f.Resource = product
	}
	return nil
}

func (f *Factory) Freeze() error {
	if err := f.Freezer.Insert(f.Resource); err != nil {
		return err
	}
	if err := f.Freezer.Process(); err != nil {
		return err
	}
	if product, err := f.Freezer.Takeout(); err != nil {
		return err
	} else {
		f.Resource = product
	}
	return nil
}

func (f *Factory) Smelt() error {
	if err := f.Smelter.Insert(f.Resource); err != nil {
		return err
	}
	if err := f.Smelter.Process(); err != nil {
		return err
	}
	if product, err := f.Smelter.Takeout(); err != nil {
		return err
	} else {
		f.Resource = product
	}
	return nil
}

func (f *Factory) Process() error {
	recipe, err := GenerateRecipe(f.From, f.To)
	if err != nil {
		return err
	}
	for _, action := range recipe {
		switch action {
		case ApplyGrinding:
			if err := f.Grind(); err != nil {
				return err
			}
		case ApplySmelting:
			if err := f.Smelt(); err != nil {
				return err
			}
		case ApplyFreezing:
			if err := f.Freeze(); err != nil {
				return err
			}
		}
	}
	return nil
}
