task :default => [:test]

task :test do
	Dir.glob('test/unit/*.rb').each do |testfile|
		ruby "#{testfile}"
	end
end

task :install do
	require_relative 'lib/version.rb'
	`gem build pac.gemspec`
	`gem install pac-#{PAC::VERSION}.gem`
end

#Example
task :rspec do
	`rspec --format='html' --out='results.html'`
end

task :coverage do
  ENV['COVERAGE'] = 'on'
  Rake::Task['test'].execute
end


