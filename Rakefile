task :default => [:test]

task :test do
	Dir.glob('test/unit/*.rb').each do |testfile|
		ruby "#{testfile}"
	end
end

task :functional_test do
	Dir.glob('test/functional/*.rb').each do |testfile|
		ruby "#{testfile} --verbose=verbose"
	end
end

task :functional_trac do
	Dir.glob('test/functional/*trac*.rb').each do |testfile|
		ruby "#{testfile} --verbose=verbose"
	end
end

task :functional_jira do
	Dir.glob('test/functional/*jira*.rb').each do |testfile|
		ruby "#{testfile} --verbose=verbose"
	end
end

task :functional_vcs do
	Dir.glob('test/functional/*vcs*.rb').each do |testfile|
		puts testfile
		ruby "#{testfile} --verbose=verbose"
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


