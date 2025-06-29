// Code generated by http://github.com/gojuno/minimock (v3.4.5). DO NOT EDIT.

package mocks

//go:generate minimock -i github.com/iamvkosarev/learning-cards/internal/module.GroupReader -o group_reader_mock.go -n GroupReaderMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/iamvkosarev/learning-cards/internal/model"
)

// GroupReaderMock implements mm_module.GroupReader
type GroupReaderMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetGroup          func(ctx context.Context, groupId model.GroupId) (gp1 *model.Group, err error)
	funcGetGroupOrigin    string
	inspectFuncGetGroup   func(ctx context.Context, groupId model.GroupId)
	afterGetGroupCounter  uint64
	beforeGetGroupCounter uint64
	GetGroupMock          mGroupReaderMockGetGroup

	funcListGroups          func(ctx context.Context, id model.UserId) (gpa1 []*model.Group, err error)
	funcListGroupsOrigin    string
	inspectFuncListGroups   func(ctx context.Context, id model.UserId)
	afterListGroupsCounter  uint64
	beforeListGroupsCounter uint64
	ListGroupsMock          mGroupReaderMockListGroups
}

// NewGroupReaderMock returns a mock for mm_module.GroupReader
func NewGroupReaderMock(t minimock.Tester) *GroupReaderMock {
	m := &GroupReaderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetGroupMock = mGroupReaderMockGetGroup{mock: m}
	m.GetGroupMock.callArgs = []*GroupReaderMockGetGroupParams{}

	m.ListGroupsMock = mGroupReaderMockListGroups{mock: m}
	m.ListGroupsMock.callArgs = []*GroupReaderMockListGroupsParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mGroupReaderMockGetGroup struct {
	optional           bool
	mock               *GroupReaderMock
	defaultExpectation *GroupReaderMockGetGroupExpectation
	expectations       []*GroupReaderMockGetGroupExpectation

	callArgs []*GroupReaderMockGetGroupParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// GroupReaderMockGetGroupExpectation specifies expectation struct of the GroupReader.GetGroup
type GroupReaderMockGetGroupExpectation struct {
	mock               *GroupReaderMock
	params             *GroupReaderMockGetGroupParams
	paramPtrs          *GroupReaderMockGetGroupParamPtrs
	expectationOrigins GroupReaderMockGetGroupExpectationOrigins
	results            *GroupReaderMockGetGroupResults
	returnOrigin       string
	Counter            uint64
}

// GroupReaderMockGetGroupParams contains parameters of the GroupReader.GetGroup
type GroupReaderMockGetGroupParams struct {
	ctx     context.Context
	groupId model.GroupId
}

// GroupReaderMockGetGroupParamPtrs contains pointers to parameters of the GroupReader.GetGroup
type GroupReaderMockGetGroupParamPtrs struct {
	ctx     *context.Context
	groupId *model.GroupId
}

// GroupReaderMockGetGroupResults contains results of the GroupReader.GetGroup
type GroupReaderMockGetGroupResults struct {
	gp1 *model.Group
	err error
}

// GroupReaderMockGetGroupOrigins contains origins of expectations of the GroupReader.GetGroup
type GroupReaderMockGetGroupExpectationOrigins struct {
	origin        string
	originCtx     string
	originGroupId string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetGroup *mGroupReaderMockGetGroup) Optional() *mGroupReaderMockGetGroup {
	mmGetGroup.optional = true
	return mmGetGroup
}

// Expect sets up expected params for GroupReader.GetGroup
func (mmGetGroup *mGroupReaderMockGetGroup) Expect(ctx context.Context, groupId model.GroupId) *mGroupReaderMockGetGroup {
	if mmGetGroup.mock.funcGetGroup != nil {
		mmGetGroup.mock.t.Fatalf("GroupReaderMock.GetGroup mock is already set by Set")
	}

	if mmGetGroup.defaultExpectation == nil {
		mmGetGroup.defaultExpectation = &GroupReaderMockGetGroupExpectation{}
	}

	if mmGetGroup.defaultExpectation.paramPtrs != nil {
		mmGetGroup.mock.t.Fatalf("GroupReaderMock.GetGroup mock is already set by ExpectParams functions")
	}

	mmGetGroup.defaultExpectation.params = &GroupReaderMockGetGroupParams{ctx, groupId}
	mmGetGroup.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmGetGroup.expectations {
		if minimock.Equal(e.params, mmGetGroup.defaultExpectation.params) {
			mmGetGroup.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetGroup.defaultExpectation.params)
		}
	}

	return mmGetGroup
}

// ExpectCtxParam1 sets up expected param ctx for GroupReader.GetGroup
func (mmGetGroup *mGroupReaderMockGetGroup) ExpectCtxParam1(ctx context.Context) *mGroupReaderMockGetGroup {
	if mmGetGroup.mock.funcGetGroup != nil {
		mmGetGroup.mock.t.Fatalf("GroupReaderMock.GetGroup mock is already set by Set")
	}

	if mmGetGroup.defaultExpectation == nil {
		mmGetGroup.defaultExpectation = &GroupReaderMockGetGroupExpectation{}
	}

	if mmGetGroup.defaultExpectation.params != nil {
		mmGetGroup.mock.t.Fatalf("GroupReaderMock.GetGroup mock is already set by Expect")
	}

	if mmGetGroup.defaultExpectation.paramPtrs == nil {
		mmGetGroup.defaultExpectation.paramPtrs = &GroupReaderMockGetGroupParamPtrs{}
	}
	mmGetGroup.defaultExpectation.paramPtrs.ctx = &ctx
	mmGetGroup.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmGetGroup
}

// ExpectGroupIdParam2 sets up expected param groupId for GroupReader.GetGroup
func (mmGetGroup *mGroupReaderMockGetGroup) ExpectGroupIdParam2(groupId model.GroupId) *mGroupReaderMockGetGroup {
	if mmGetGroup.mock.funcGetGroup != nil {
		mmGetGroup.mock.t.Fatalf("GroupReaderMock.GetGroup mock is already set by Set")
	}

	if mmGetGroup.defaultExpectation == nil {
		mmGetGroup.defaultExpectation = &GroupReaderMockGetGroupExpectation{}
	}

	if mmGetGroup.defaultExpectation.params != nil {
		mmGetGroup.mock.t.Fatalf("GroupReaderMock.GetGroup mock is already set by Expect")
	}

	if mmGetGroup.defaultExpectation.paramPtrs == nil {
		mmGetGroup.defaultExpectation.paramPtrs = &GroupReaderMockGetGroupParamPtrs{}
	}
	mmGetGroup.defaultExpectation.paramPtrs.groupId = &groupId
	mmGetGroup.defaultExpectation.expectationOrigins.originGroupId = minimock.CallerInfo(1)

	return mmGetGroup
}

// Inspect accepts an inspector function that has same arguments as the GroupReader.GetGroup
func (mmGetGroup *mGroupReaderMockGetGroup) Inspect(f func(ctx context.Context, groupId model.GroupId)) *mGroupReaderMockGetGroup {
	if mmGetGroup.mock.inspectFuncGetGroup != nil {
		mmGetGroup.mock.t.Fatalf("Inspect function is already set for GroupReaderMock.GetGroup")
	}

	mmGetGroup.mock.inspectFuncGetGroup = f

	return mmGetGroup
}

// Return sets up results that will be returned by GroupReader.GetGroup
func (mmGetGroup *mGroupReaderMockGetGroup) Return(gp1 *model.Group, err error) *GroupReaderMock {
	if mmGetGroup.mock.funcGetGroup != nil {
		mmGetGroup.mock.t.Fatalf("GroupReaderMock.GetGroup mock is already set by Set")
	}

	if mmGetGroup.defaultExpectation == nil {
		mmGetGroup.defaultExpectation = &GroupReaderMockGetGroupExpectation{mock: mmGetGroup.mock}
	}
	mmGetGroup.defaultExpectation.results = &GroupReaderMockGetGroupResults{gp1, err}
	mmGetGroup.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmGetGroup.mock
}

// Set uses given function f to mock the GroupReader.GetGroup method
func (mmGetGroup *mGroupReaderMockGetGroup) Set(f func(ctx context.Context, groupId model.GroupId) (gp1 *model.Group, err error)) *GroupReaderMock {
	if mmGetGroup.defaultExpectation != nil {
		mmGetGroup.mock.t.Fatalf("Default expectation is already set for the GroupReader.GetGroup method")
	}

	if len(mmGetGroup.expectations) > 0 {
		mmGetGroup.mock.t.Fatalf("Some expectations are already set for the GroupReader.GetGroup method")
	}

	mmGetGroup.mock.funcGetGroup = f
	mmGetGroup.mock.funcGetGroupOrigin = minimock.CallerInfo(1)
	return mmGetGroup.mock
}

// When sets expectation for the GroupReader.GetGroup which will trigger the result defined by the following
// Then helper
func (mmGetGroup *mGroupReaderMockGetGroup) When(ctx context.Context, groupId model.GroupId) *GroupReaderMockGetGroupExpectation {
	if mmGetGroup.mock.funcGetGroup != nil {
		mmGetGroup.mock.t.Fatalf("GroupReaderMock.GetGroup mock is already set by Set")
	}

	expectation := &GroupReaderMockGetGroupExpectation{
		mock:               mmGetGroup.mock,
		params:             &GroupReaderMockGetGroupParams{ctx, groupId},
		expectationOrigins: GroupReaderMockGetGroupExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmGetGroup.expectations = append(mmGetGroup.expectations, expectation)
	return expectation
}

// Then sets up GroupReader.GetGroup return parameters for the expectation previously defined by the When method
func (e *GroupReaderMockGetGroupExpectation) Then(gp1 *model.Group, err error) *GroupReaderMock {
	e.results = &GroupReaderMockGetGroupResults{gp1, err}
	return e.mock
}

// Times sets number of times GroupReader.GetGroup should be invoked
func (mmGetGroup *mGroupReaderMockGetGroup) Times(n uint64) *mGroupReaderMockGetGroup {
	if n == 0 {
		mmGetGroup.mock.t.Fatalf("Times of GroupReaderMock.GetGroup mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetGroup.expectedInvocations, n)
	mmGetGroup.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmGetGroup
}

func (mmGetGroup *mGroupReaderMockGetGroup) invocationsDone() bool {
	if len(mmGetGroup.expectations) == 0 && mmGetGroup.defaultExpectation == nil && mmGetGroup.mock.funcGetGroup == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetGroup.mock.afterGetGroupCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetGroup.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetGroup implements mm_module.GroupReader
func (mmGetGroup *GroupReaderMock) GetGroup(ctx context.Context, groupId model.GroupId) (gp1 *model.Group, err error) {
	mm_atomic.AddUint64(&mmGetGroup.beforeGetGroupCounter, 1)
	defer mm_atomic.AddUint64(&mmGetGroup.afterGetGroupCounter, 1)

	mmGetGroup.t.Helper()

	if mmGetGroup.inspectFuncGetGroup != nil {
		mmGetGroup.inspectFuncGetGroup(ctx, groupId)
	}

	mm_params := GroupReaderMockGetGroupParams{ctx, groupId}

	// Record call args
	mmGetGroup.GetGroupMock.mutex.Lock()
	mmGetGroup.GetGroupMock.callArgs = append(mmGetGroup.GetGroupMock.callArgs, &mm_params)
	mmGetGroup.GetGroupMock.mutex.Unlock()

	for _, e := range mmGetGroup.GetGroupMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.gp1, e.results.err
		}
	}

	if mmGetGroup.GetGroupMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetGroup.GetGroupMock.defaultExpectation.Counter, 1)
		mm_want := mmGetGroup.GetGroupMock.defaultExpectation.params
		mm_want_ptrs := mmGetGroup.GetGroupMock.defaultExpectation.paramPtrs

		mm_got := GroupReaderMockGetGroupParams{ctx, groupId}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmGetGroup.t.Errorf("GroupReaderMock.GetGroup got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmGetGroup.GetGroupMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.groupId != nil && !minimock.Equal(*mm_want_ptrs.groupId, mm_got.groupId) {
				mmGetGroup.t.Errorf("GroupReaderMock.GetGroup got unexpected parameter groupId, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmGetGroup.GetGroupMock.defaultExpectation.expectationOrigins.originGroupId, *mm_want_ptrs.groupId, mm_got.groupId, minimock.Diff(*mm_want_ptrs.groupId, mm_got.groupId))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetGroup.t.Errorf("GroupReaderMock.GetGroup got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmGetGroup.GetGroupMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetGroup.GetGroupMock.defaultExpectation.results
		if mm_results == nil {
			mmGetGroup.t.Fatal("No results are set for the GroupReaderMock.GetGroup")
		}
		return (*mm_results).gp1, (*mm_results).err
	}
	if mmGetGroup.funcGetGroup != nil {
		return mmGetGroup.funcGetGroup(ctx, groupId)
	}
	mmGetGroup.t.Fatalf("Unexpected call to GroupReaderMock.GetGroup. %v %v", ctx, groupId)
	return
}

// GetGroupAfterCounter returns a count of finished GroupReaderMock.GetGroup invocations
func (mmGetGroup *GroupReaderMock) GetGroupAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetGroup.afterGetGroupCounter)
}

// GetGroupBeforeCounter returns a count of GroupReaderMock.GetGroup invocations
func (mmGetGroup *GroupReaderMock) GetGroupBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetGroup.beforeGetGroupCounter)
}

// Calls returns a list of arguments used in each call to GroupReaderMock.GetGroup.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetGroup *mGroupReaderMockGetGroup) Calls() []*GroupReaderMockGetGroupParams {
	mmGetGroup.mutex.RLock()

	argCopy := make([]*GroupReaderMockGetGroupParams, len(mmGetGroup.callArgs))
	copy(argCopy, mmGetGroup.callArgs)

	mmGetGroup.mutex.RUnlock()

	return argCopy
}

// MinimockGetGroupDone returns true if the count of the GetGroup invocations corresponds
// the number of defined expectations
func (m *GroupReaderMock) MinimockGetGroupDone() bool {
	if m.GetGroupMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetGroupMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetGroupMock.invocationsDone()
}

// MinimockGetGroupInspect logs each unmet expectation
func (m *GroupReaderMock) MinimockGetGroupInspect() {
	for _, e := range m.GetGroupMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to GroupReaderMock.GetGroup at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterGetGroupCounter := mm_atomic.LoadUint64(&m.afterGetGroupCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetGroupMock.defaultExpectation != nil && afterGetGroupCounter < 1 {
		if m.GetGroupMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to GroupReaderMock.GetGroup at\n%s", m.GetGroupMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to GroupReaderMock.GetGroup at\n%s with params: %#v", m.GetGroupMock.defaultExpectation.expectationOrigins.origin, *m.GetGroupMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetGroup != nil && afterGetGroupCounter < 1 {
		m.t.Errorf("Expected call to GroupReaderMock.GetGroup at\n%s", m.funcGetGroupOrigin)
	}

	if !m.GetGroupMock.invocationsDone() && afterGetGroupCounter > 0 {
		m.t.Errorf("Expected %d calls to GroupReaderMock.GetGroup at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.GetGroupMock.expectedInvocations), m.GetGroupMock.expectedInvocationsOrigin, afterGetGroupCounter)
	}
}

type mGroupReaderMockListGroups struct {
	optional           bool
	mock               *GroupReaderMock
	defaultExpectation *GroupReaderMockListGroupsExpectation
	expectations       []*GroupReaderMockListGroupsExpectation

	callArgs []*GroupReaderMockListGroupsParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// GroupReaderMockListGroupsExpectation specifies expectation struct of the GroupReader.ListGroups
type GroupReaderMockListGroupsExpectation struct {
	mock               *GroupReaderMock
	params             *GroupReaderMockListGroupsParams
	paramPtrs          *GroupReaderMockListGroupsParamPtrs
	expectationOrigins GroupReaderMockListGroupsExpectationOrigins
	results            *GroupReaderMockListGroupsResults
	returnOrigin       string
	Counter            uint64
}

// GroupReaderMockListGroupsParams contains parameters of the GroupReader.ListGroups
type GroupReaderMockListGroupsParams struct {
	ctx context.Context
	id  model.UserId
}

// GroupReaderMockListGroupsParamPtrs contains pointers to parameters of the GroupReader.ListGroups
type GroupReaderMockListGroupsParamPtrs struct {
	ctx *context.Context
	id  *model.UserId
}

// GroupReaderMockListGroupsResults contains results of the GroupReader.ListGroups
type GroupReaderMockListGroupsResults struct {
	gpa1 []*model.Group
	err  error
}

// GroupReaderMockListGroupsOrigins contains origins of expectations of the GroupReader.ListGroups
type GroupReaderMockListGroupsExpectationOrigins struct {
	origin    string
	originCtx string
	originId  string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmListGroups *mGroupReaderMockListGroups) Optional() *mGroupReaderMockListGroups {
	mmListGroups.optional = true
	return mmListGroups
}

// Expect sets up expected params for GroupReader.ListGroups
func (mmListGroups *mGroupReaderMockListGroups) Expect(ctx context.Context, id model.UserId) *mGroupReaderMockListGroups {
	if mmListGroups.mock.funcListGroups != nil {
		mmListGroups.mock.t.Fatalf("GroupReaderMock.ListGroups mock is already set by Set")
	}

	if mmListGroups.defaultExpectation == nil {
		mmListGroups.defaultExpectation = &GroupReaderMockListGroupsExpectation{}
	}

	if mmListGroups.defaultExpectation.paramPtrs != nil {
		mmListGroups.mock.t.Fatalf("GroupReaderMock.ListGroups mock is already set by ExpectParams functions")
	}

	mmListGroups.defaultExpectation.params = &GroupReaderMockListGroupsParams{ctx, id}
	mmListGroups.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmListGroups.expectations {
		if minimock.Equal(e.params, mmListGroups.defaultExpectation.params) {
			mmListGroups.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmListGroups.defaultExpectation.params)
		}
	}

	return mmListGroups
}

// ExpectCtxParam1 sets up expected param ctx for GroupReader.ListGroups
func (mmListGroups *mGroupReaderMockListGroups) ExpectCtxParam1(ctx context.Context) *mGroupReaderMockListGroups {
	if mmListGroups.mock.funcListGroups != nil {
		mmListGroups.mock.t.Fatalf("GroupReaderMock.ListGroups mock is already set by Set")
	}

	if mmListGroups.defaultExpectation == nil {
		mmListGroups.defaultExpectation = &GroupReaderMockListGroupsExpectation{}
	}

	if mmListGroups.defaultExpectation.params != nil {
		mmListGroups.mock.t.Fatalf("GroupReaderMock.ListGroups mock is already set by Expect")
	}

	if mmListGroups.defaultExpectation.paramPtrs == nil {
		mmListGroups.defaultExpectation.paramPtrs = &GroupReaderMockListGroupsParamPtrs{}
	}
	mmListGroups.defaultExpectation.paramPtrs.ctx = &ctx
	mmListGroups.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmListGroups
}

// ExpectIdParam2 sets up expected param id for GroupReader.ListGroups
func (mmListGroups *mGroupReaderMockListGroups) ExpectIdParam2(id model.UserId) *mGroupReaderMockListGroups {
	if mmListGroups.mock.funcListGroups != nil {
		mmListGroups.mock.t.Fatalf("GroupReaderMock.ListGroups mock is already set by Set")
	}

	if mmListGroups.defaultExpectation == nil {
		mmListGroups.defaultExpectation = &GroupReaderMockListGroupsExpectation{}
	}

	if mmListGroups.defaultExpectation.params != nil {
		mmListGroups.mock.t.Fatalf("GroupReaderMock.ListGroups mock is already set by Expect")
	}

	if mmListGroups.defaultExpectation.paramPtrs == nil {
		mmListGroups.defaultExpectation.paramPtrs = &GroupReaderMockListGroupsParamPtrs{}
	}
	mmListGroups.defaultExpectation.paramPtrs.id = &id
	mmListGroups.defaultExpectation.expectationOrigins.originId = minimock.CallerInfo(1)

	return mmListGroups
}

// Inspect accepts an inspector function that has same arguments as the GroupReader.ListGroups
func (mmListGroups *mGroupReaderMockListGroups) Inspect(f func(ctx context.Context, id model.UserId)) *mGroupReaderMockListGroups {
	if mmListGroups.mock.inspectFuncListGroups != nil {
		mmListGroups.mock.t.Fatalf("Inspect function is already set for GroupReaderMock.ListGroups")
	}

	mmListGroups.mock.inspectFuncListGroups = f

	return mmListGroups
}

// Return sets up results that will be returned by GroupReader.ListGroups
func (mmListGroups *mGroupReaderMockListGroups) Return(gpa1 []*model.Group, err error) *GroupReaderMock {
	if mmListGroups.mock.funcListGroups != nil {
		mmListGroups.mock.t.Fatalf("GroupReaderMock.ListGroups mock is already set by Set")
	}

	if mmListGroups.defaultExpectation == nil {
		mmListGroups.defaultExpectation = &GroupReaderMockListGroupsExpectation{mock: mmListGroups.mock}
	}
	mmListGroups.defaultExpectation.results = &GroupReaderMockListGroupsResults{gpa1, err}
	mmListGroups.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmListGroups.mock
}

// Set uses given function f to mock the GroupReader.ListGroups method
func (mmListGroups *mGroupReaderMockListGroups) Set(f func(ctx context.Context, id model.UserId) (gpa1 []*model.Group, err error)) *GroupReaderMock {
	if mmListGroups.defaultExpectation != nil {
		mmListGroups.mock.t.Fatalf("Default expectation is already set for the GroupReader.ListGroups method")
	}

	if len(mmListGroups.expectations) > 0 {
		mmListGroups.mock.t.Fatalf("Some expectations are already set for the GroupReader.ListGroups method")
	}

	mmListGroups.mock.funcListGroups = f
	mmListGroups.mock.funcListGroupsOrigin = minimock.CallerInfo(1)
	return mmListGroups.mock
}

// When sets expectation for the GroupReader.ListGroups which will trigger the result defined by the following
// Then helper
func (mmListGroups *mGroupReaderMockListGroups) When(ctx context.Context, id model.UserId) *GroupReaderMockListGroupsExpectation {
	if mmListGroups.mock.funcListGroups != nil {
		mmListGroups.mock.t.Fatalf("GroupReaderMock.ListGroups mock is already set by Set")
	}

	expectation := &GroupReaderMockListGroupsExpectation{
		mock:               mmListGroups.mock,
		params:             &GroupReaderMockListGroupsParams{ctx, id},
		expectationOrigins: GroupReaderMockListGroupsExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmListGroups.expectations = append(mmListGroups.expectations, expectation)
	return expectation
}

// Then sets up GroupReader.ListGroups return parameters for the expectation previously defined by the When method
func (e *GroupReaderMockListGroupsExpectation) Then(gpa1 []*model.Group, err error) *GroupReaderMock {
	e.results = &GroupReaderMockListGroupsResults{gpa1, err}
	return e.mock
}

// Times sets number of times GroupReader.ListGroups should be invoked
func (mmListGroups *mGroupReaderMockListGroups) Times(n uint64) *mGroupReaderMockListGroups {
	if n == 0 {
		mmListGroups.mock.t.Fatalf("Times of GroupReaderMock.ListGroups mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmListGroups.expectedInvocations, n)
	mmListGroups.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmListGroups
}

func (mmListGroups *mGroupReaderMockListGroups) invocationsDone() bool {
	if len(mmListGroups.expectations) == 0 && mmListGroups.defaultExpectation == nil && mmListGroups.mock.funcListGroups == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmListGroups.mock.afterListGroupsCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmListGroups.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// ListGroups implements mm_module.GroupReader
func (mmListGroups *GroupReaderMock) ListGroups(ctx context.Context, id model.UserId) (gpa1 []*model.Group, err error) {
	mm_atomic.AddUint64(&mmListGroups.beforeListGroupsCounter, 1)
	defer mm_atomic.AddUint64(&mmListGroups.afterListGroupsCounter, 1)

	mmListGroups.t.Helper()

	if mmListGroups.inspectFuncListGroups != nil {
		mmListGroups.inspectFuncListGroups(ctx, id)
	}

	mm_params := GroupReaderMockListGroupsParams{ctx, id}

	// Record call args
	mmListGroups.ListGroupsMock.mutex.Lock()
	mmListGroups.ListGroupsMock.callArgs = append(mmListGroups.ListGroupsMock.callArgs, &mm_params)
	mmListGroups.ListGroupsMock.mutex.Unlock()

	for _, e := range mmListGroups.ListGroupsMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.gpa1, e.results.err
		}
	}

	if mmListGroups.ListGroupsMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmListGroups.ListGroupsMock.defaultExpectation.Counter, 1)
		mm_want := mmListGroups.ListGroupsMock.defaultExpectation.params
		mm_want_ptrs := mmListGroups.ListGroupsMock.defaultExpectation.paramPtrs

		mm_got := GroupReaderMockListGroupsParams{ctx, id}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmListGroups.t.Errorf("GroupReaderMock.ListGroups got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmListGroups.ListGroupsMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.id != nil && !minimock.Equal(*mm_want_ptrs.id, mm_got.id) {
				mmListGroups.t.Errorf("GroupReaderMock.ListGroups got unexpected parameter id, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmListGroups.ListGroupsMock.defaultExpectation.expectationOrigins.originId, *mm_want_ptrs.id, mm_got.id, minimock.Diff(*mm_want_ptrs.id, mm_got.id))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmListGroups.t.Errorf("GroupReaderMock.ListGroups got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmListGroups.ListGroupsMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmListGroups.ListGroupsMock.defaultExpectation.results
		if mm_results == nil {
			mmListGroups.t.Fatal("No results are set for the GroupReaderMock.ListGroups")
		}
		return (*mm_results).gpa1, (*mm_results).err
	}
	if mmListGroups.funcListGroups != nil {
		return mmListGroups.funcListGroups(ctx, id)
	}
	mmListGroups.t.Fatalf("Unexpected call to GroupReaderMock.ListGroups. %v %v", ctx, id)
	return
}

// ListGroupsAfterCounter returns a count of finished GroupReaderMock.ListGroups invocations
func (mmListGroups *GroupReaderMock) ListGroupsAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmListGroups.afterListGroupsCounter)
}

// ListGroupsBeforeCounter returns a count of GroupReaderMock.ListGroups invocations
func (mmListGroups *GroupReaderMock) ListGroupsBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmListGroups.beforeListGroupsCounter)
}

// Calls returns a list of arguments used in each call to GroupReaderMock.ListGroups.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmListGroups *mGroupReaderMockListGroups) Calls() []*GroupReaderMockListGroupsParams {
	mmListGroups.mutex.RLock()

	argCopy := make([]*GroupReaderMockListGroupsParams, len(mmListGroups.callArgs))
	copy(argCopy, mmListGroups.callArgs)

	mmListGroups.mutex.RUnlock()

	return argCopy
}

// MinimockListGroupsDone returns true if the count of the ListGroups invocations corresponds
// the number of defined expectations
func (m *GroupReaderMock) MinimockListGroupsDone() bool {
	if m.ListGroupsMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.ListGroupsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.ListGroupsMock.invocationsDone()
}

// MinimockListGroupsInspect logs each unmet expectation
func (m *GroupReaderMock) MinimockListGroupsInspect() {
	for _, e := range m.ListGroupsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to GroupReaderMock.ListGroups at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterListGroupsCounter := mm_atomic.LoadUint64(&m.afterListGroupsCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.ListGroupsMock.defaultExpectation != nil && afterListGroupsCounter < 1 {
		if m.ListGroupsMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to GroupReaderMock.ListGroups at\n%s", m.ListGroupsMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to GroupReaderMock.ListGroups at\n%s with params: %#v", m.ListGroupsMock.defaultExpectation.expectationOrigins.origin, *m.ListGroupsMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcListGroups != nil && afterListGroupsCounter < 1 {
		m.t.Errorf("Expected call to GroupReaderMock.ListGroups at\n%s", m.funcListGroupsOrigin)
	}

	if !m.ListGroupsMock.invocationsDone() && afterListGroupsCounter > 0 {
		m.t.Errorf("Expected %d calls to GroupReaderMock.ListGroups at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.ListGroupsMock.expectedInvocations), m.ListGroupsMock.expectedInvocationsOrigin, afterListGroupsCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *GroupReaderMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetGroupInspect()

			m.MinimockListGroupsInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *GroupReaderMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *GroupReaderMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetGroupDone() &&
		m.MinimockListGroupsDone()
}
