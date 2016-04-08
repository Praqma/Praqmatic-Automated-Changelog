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

task :coverage do
  ENV['COVERAGE'] = 'on'
  Rake::Task['test'].execute
end


