package temporaltestcancel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(string, string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) TestWorkflow() {
	var timerCancelled bool
	s.env.SetOnTimerCanceledListener(func(timerID string) {
		timerCancelled = true
	})

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("FireTimer", nil)
	}, 24*time.Hour)

	s.env.ExecuteWorkflow(Workflow)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
	s.True(timerCancelled)
}

func TestWorkflow(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
