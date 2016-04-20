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

task :coverage do
  ENV['COVERAGE'] = 'on'
  Rake::Task['test'].execute
end


