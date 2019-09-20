# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'version'

Gem::Specification.new do |spec|
  spec.name          = "pac"
  spec.version       = PAC::VERSION
  spec.authors       = ["Mads Nielsen"]
  spec.email         = ["man@praqma.net"]

  spec.summary       = %q{Praqmatic Automated Changelog}
  spec.description   = %q{Use this gem to create a release note from git repo and task system}
  spec.homepage      = "https://github.com/Praqma/Praqmatic-Automated-Changelog"
  spec.license       = "MIT"

  spec.files         = `git ls-files -z`.split("\x0").reject do |f|
    f.match(%r{^(test|spec|features|jenkins-pipeline|site|templates|settings|reports)/})
  end

  spec.bindir        = "bin"
  spec.executables   = 'pac'
  spec.require_paths = ["lib"]

  spec.add_development_dependency 'bundler', '~> 1.14'
  spec.add_development_dependency 'rake', '~> 10.0'
  spec.add_development_dependency 'rspec', '~> 3.7'
  spec.add_development_dependency 'test-unit', '~> 2.5.5', '>= 2.5.5'

  spec.add_runtime_dependency 'docopt', '~> 0.6.1', '>= 0.6.1'
  spec.add_runtime_dependency 'rugged', '~> 0.26', '>= 0.26'
  spec.add_runtime_dependency 'liquid', '~> 4.0.0', '>= 4.0.0'

end
