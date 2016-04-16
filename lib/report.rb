module Report

	class Generator
	  #This one has become a tad too specific for git. We get the short shar (first 8 digits). And
	  #we trim the subject (That is the first element in the list when splitting the message on a newline. 
	  #Arguments:
	  # tasks     - A 'PACTaskCollection' of populated tasks
	  # commits   - A 'PACCommitCollection' list of commits. used to tally those that were referenced and those that were not
	  # templates - The list of configured templates to render 
	  def generate(tasks, commits, templates)
	    templates.each do |t|        
	      File.open(t['output'],'w:UTF-8') do |file| 
	        file << Liquid::Template.parse(File.read(t['location'])).render( { 
	          'tasks' => tasks, 	     
	          'title' => 'PAC id report',
	          'pac_c_count' => commits.count,
	          'pac_c_referenced' => commits.count_with,
	          'pac_health' => commits.health,
	          'pac_c_unreferenced' => commits.count_without  
	        } )
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
	end
end