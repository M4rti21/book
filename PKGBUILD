# Maintainer: Martí Comas <m4rti21@proton.me>
pkgname='book'
pkgver=0.2.4
pkgrel=1
pkgdesc="A simple bookmark manager written in go"
arch=('x86_64')
url="https://github.com/M4rti21/$pkgname"
license=('GPL')
groups=()
depends=()
makedepends=('go')
optdepends=()
source=("$pkgname.tar.gz::https://github.com/M4rti21/$pkgname/archive/refs/tags/$pkgver.tar.gz")
sha256sums=('SKIP')
options=('!debug')

build() {
    cd "$pkgname-$pkgver/src"
    go build -o ..
}

package() {
    cd "$pkgname-$pkgver"
    install -Dm755 "./$pkgname" "$pkgdir/usr/bin/$pkgname"
}
