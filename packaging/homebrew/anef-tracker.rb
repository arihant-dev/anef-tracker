class AnefTracker < Formula
  desc "Private, Multi-Application Evidence Vault & Workflow Intelligence Platform for ANEF"
  homepage "https://github.com/anef-tracker/anef-tracker"
  url "https://github.com/anef-tracker/anef-tracker/releases/download/v0.9.0/anef-tracker_0.9.0_darwin_arm64.tar.gz"
  sha256 "0000000000000000000000000000000000000000000000000000000000000000"
  license "MIT"

  def install
    bin.install "anef"
  end

  test do
    assert_match "ANEF Tracker", shell_output("#{bin}/anef version")
  end
end
