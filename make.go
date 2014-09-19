package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const gtcproxyRbFmt = `# WARNING: Automatically generated. All changes to this file will be lost
require 'formula'

HOMEBREW_PECO_VERSION='%s'
class Peco < Formula
  homepage 'https://github.com/lestrrat/gtcproxy'
  url "https://github.com/lestrrat/gtcproxy/releases/download/v#{HOMEBREW_PECO_VERSION}/gtcproxy_darwin_amd64.zip"
  sha1 "%x"

  version HOMEBREW_PECO_VERSION
  head 'https://github.com/lestrrat/gtcproxy.git', :branch => 'master'

  if build.head?
    depends_on 'go' => :build
    depends_on 'hg' => :build
  end

  def install
    if build.head?
      ENV['GOPATH'] = buildpath
      system 'go', 'build', '.'
    end
    bin.install 'gtcproxy'
  end
end`

func main() {
	st := _main()
	os.Exit(st)
}

func _main() int {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage:\n  go run make.go [version]\n")
		return 1
	}

	version := os.Args[1]
	return updateGenericRb("gtcproxy", version, gtcproxyRbFmt)
}

// fetch applicable binary, generate checksum, and update the .rb file
func updateGenericRb(target, version, template string) int {
	url := fmt.Sprintf(
		"https://github.com/lestrrat/%s/releases/download/v%s/%s_darwin_amd64.zip",
		target,
		version,
		target,
	)

	log.Printf("Fetching url %s...", url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if res.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Got %d instead of 200", res.StatusCode)
		return 1
	}

	h := sha1.New()
	io.Copy(h, res.Body)

	filename := fmt.Sprintf("%s.rb", target)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file %s: %s", filename, err)
		return 1
	}

	fmt.Fprintf(file, template, version, h.Sum(nil))
	return 0
}
