# WARNING: Automatically generated. All changes to this file will be lost
require 'formula'

HOMEBREW_GTCPROXY_VERSION='0.0.1'
class Gtcproxy < Formula
  homepage 'https://github.com/lestrrat/gtcproxy'
  url "https://github.com/lestrrat/gtcproxy/releases/download/v#{HOMEBREW_PECO_VERSION}/gtcproxy_darwin_amd64.zip"
  sha1 "78f2e67425c263fdfadab811834e10e7b6994602"

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
end