package app

type Service interface {
    GetAssets() ([]Asset, error)
    GetAsset(name string) (Asset, error)
}

type Asset struct {
    Name string
    Symbol string
}
