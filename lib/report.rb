module Report

	class Generator
	  #This one has become a tad too specific for git. We get the short shar (first 8 digits). And
	  #we trim the subject (That is the first element in the list when splitting the message on a newline. 
	  #Arguments:
	  # tasks     - A 'PACTaskCollection' of populated tasks
	  # commits   - A 'PACCommitCollection' list of commits. used to tally those that were referenced and those that were not
	  # config 		- The configuration used 
	  def generate(tasks, commits, config)
	    config[:templates].each do |t|        
	      File.open(t['output'],'w:UTF-8') do |file| 
	        file << Liquid::Template.parse(File.read(t['location'])).render( { 
	          'tasks' => tasks, 	     
	          'pac_c_count' => commits.count,
	          'pac_c_referenced' => commits.count_with,
	          'pac_health' => commits.health,
	          'pac_c_unreferenced' => commits.count_without  
	        }.merge!(define_properties(config)) )
	        File.chmod(0777, file)
	      end

	      if t['pdf'] == true
	        output_file_path = t['output'].sub(/\.html$/, '.pdf')
	        kit = PDFKit.new(File.new(t['output']), :page_size => 'A4')
	        kit.to_file(output_file_path)
	        File.chmod(0777, output_file_path)
	      end
	    end
	  end

	  #Merges the properties, defaults title to 'PAC Changelog'
	  #At this point config[:properties] is either nil, or properly configured and defined.
	  def define_properties(config)
	  	std = { 'properties' => { 'title' => 'PAC Changelog' } }
	  	unless config[:properties].nil?
	  		std['properties'].merge!(config[:properties])
	  	end
	  	std
	  end 
	end
end