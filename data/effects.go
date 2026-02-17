package data

type Effect_WaterColor uint32

func (*Effect_WaterColor) ID() string { return "water_color" }

type Effect_FoliageColor uint32

func (*Effect_FoliageColor) ID() string { return "foliage_color" }

type Effect_DryFoliageColor uint32

func (*Effect_DryFoliageColor) ID() string { return "dry_foliage_color" }

type Effect_GrassColor uint32

func (*Effect_GrassColor) ID() string { return "grass_color" }

type Effect_GrassColorModifier string

func (*Effect_GrassColorModifier) ID() string { return "grass_color_modifier" }
