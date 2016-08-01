# In general the test cases here verifies we collect the correct commits
# on the relevant branch, and nothing more.
module PAC__TestSuites
    # TestCase classes that contain the methods I want to use in the test suite
    require_relative 'testCases/vcs_functional_walker.rb'
    require_relative 'testCases/vcs_functional_topology.rb'
    include PAC__TestCases_GetCommitMessageOnCorrectBranch
end