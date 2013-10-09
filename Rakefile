task :default => [:test]

task :test do
  ruby 'tests/unit/vcstests.rb'
end

task :coverage do
  ENV['COVERAGE'] = 'on'
  Rake::Task['test'].execute
end
