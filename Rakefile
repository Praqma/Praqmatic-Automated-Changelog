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

task :changelog do
	`docker run --rm -v $(pwd):/data -v /home/praqma/tools:/tools praqma/pac:2.1.0-beta --from-latest-tag 2.0.1 --settings=/tools/config/settings_github_pac.yml -vvv` 
end

task :coverage do
  ENV['COVERAGE'] = 'on'
  Rake::Task['test'].execute
end


