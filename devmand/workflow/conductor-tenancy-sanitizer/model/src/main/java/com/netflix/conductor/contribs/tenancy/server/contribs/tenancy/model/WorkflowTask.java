package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model;

import java.util.HashMap;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Set;
import javax.validation.Valid;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.PositiveOrZero;

/**
 * @author Viren
 *
 * This is the task definition definied as part of the {@link WorkflowDef}. The tasks definied in the Workflow definition are saved
 * as part of {@link WorkflowDef#getTasks}
 */
public class WorkflowTask {

	/**
	 * This field is deprecated and will be removed in the next version.
	 * Please use {@link TaskType} instead.
	 */
	@Deprecated
	public enum Type {
		SIMPLE, DYNAMIC, FORK_JOIN, FORK_JOIN_DYNAMIC, DECISION, JOIN, SUB_WORKFLOW, EVENT, WAIT, USER_DEFINED;
		private static Set<String> systemTasks = new HashSet<>();
		static {
			systemTasks.add(Type.SIMPLE.name());
			systemTasks.add(Type.DYNAMIC.name());
			systemTasks.add(Type.FORK_JOIN.name());
			systemTasks.add(Type.FORK_JOIN_DYNAMIC.name());
			systemTasks.add(Type.DECISION.name());
			systemTasks.add(Type.JOIN.name());
			systemTasks.add(Type.SUB_WORKFLOW.name());
			systemTasks.add(Type.EVENT.name());
			systemTasks.add(Type.WAIT.name());
			//Do NOT add USER_DEFINED here...
		}
		public static boolean isSystemTask(String name) {
			return systemTasks.contains(name);
		}
	}


	@NotEmpty(message = "WorkflowTask name cannot be empty or null")
	private String name;


	@NotEmpty(message = "WorkflowTask taskReferenceName name cannot be empty or null")
	private String taskReferenceName;


	private String description;


	private Map<String, Object> inputParameters = new HashMap<>();


	private String type = TaskType.SIMPLE.name();


	private String dynamicTaskNameParam;


	private String caseValueParam;


	private String caseExpression;


	private String scriptExpression;


	public static class WorkflowTaskList {
		public List<WorkflowTask> getTasks() {
			return tasks;
		}

		public void setTasks(List<WorkflowTask> tasks) {
			this.tasks = tasks;
		}


		private List<WorkflowTask> tasks;
	}

	//Populates for the tasks of the decision type

	private Map<String,@Valid List<@Valid WorkflowTask>> decisionCases = new LinkedHashMap<>();

	@Deprecated
	private String dynamicForkJoinTasksParam;


	private String dynamicForkTasksParam;


	private String dynamicForkTasksInputParamName;


	private List<@Valid WorkflowTask> defaultCase = new LinkedList<>();


	private List<@Valid List<@Valid WorkflowTask>> forkTasks = new LinkedList<>();


    @PositiveOrZero
	private int startDelay;	//No. of seconds (at-least) to wait before starting a task.


    @Valid
	private SubWorkflowParams subWorkflowParam;


	private List<String> joinOn = new LinkedList<>();


	private String sink;


	private boolean optional = false;


	private TaskDef taskDefinition;


	private Boolean rateLimited;


	private List<String> defaultExclusiveJoinTask = new LinkedList<>();


	private Boolean asyncComplete = false;


	private String loopCondition;


	private List<WorkflowTask> loopOver = new LinkedList<>();

	/**
	 * @return the name
	 */
	public String getName() {
		return name;
	}

	/**
	 * @param name the name to set
	 */
	public void setName(String name) {
		this.name = name;
	}

	/**
	 * @return the taskReferenceName
	 */
	public String getTaskReferenceName() {
		return taskReferenceName;
	}

	/**
	 * @param taskReferenceName the taskReferenceName to set
	 */
	public void setTaskReferenceName(String taskReferenceName) {
		this.taskReferenceName = taskReferenceName;
	}

	/**
	 * @return the description
	 */
	public String getDescription() {
		return description;
	}

	/**
	 * @param description the description to set
	 */
	public void setDescription(String description) {
		this.description = description;
	}

	/**
	 * @return the inputParameters
	 */
	public Map<String, Object> getInputParameters() {
		return inputParameters;
	}

	/**
	 * @param inputParameters the inputParameters to set
	 */
	public void setInputParameters(Map<String, Object> inputParameters) {
		this.inputParameters = inputParameters;
	}

	/**
	 * @return the type
	 */
	public String getType() {
		return type;
	}

	public void setWorkflowTaskType(TaskType type) {
		this.type = type.name();
	}

	/**
	 * @param type the type to set
	 */
	public void setType(@NotEmpty(message = "WorkTask type cannot be null or empty") String type) {
		this.type = type;
	}

	/**
	 * @return the decisionCases
	 */
	public Map<String, List<WorkflowTask>> getDecisionCases() {
		return decisionCases;
	}

	/**
	 * @param decisionCases the decisionCases to set
	 */
	public void setDecisionCases(Map<String, List<WorkflowTask>> decisionCases) {
		this.decisionCases = decisionCases;
	}

	/**
	 * @return the defaultCase
	 */
	public List<WorkflowTask> getDefaultCase() {
		return defaultCase;
	}

	/**
	 * @param defaultCase the defaultCase to set
	 */
	public void setDefaultCase(List<WorkflowTask> defaultCase) {
		this.defaultCase = defaultCase;
	}

	/**
	 * @return the forkTasks
	 */
	public List<List<WorkflowTask>> getForkTasks() {
		return forkTasks;
	}

	/**
	 * @param forkTasks the forkTasks to set
	 */
	public void setForkTasks(List<List<WorkflowTask>> forkTasks) {
		this.forkTasks = forkTasks;
	}

	/**
	 * @return the startDelay in seconds
	 */
	public int getStartDelay() {
		return startDelay;
	}

	/**
	 * @param startDelay the startDelay to set
	 */
	public void setStartDelay(int startDelay) {
		this.startDelay = startDelay;
	}


	/**
	 * @return the dynamicTaskNameParam
	 */
	public String getDynamicTaskNameParam() {
		return dynamicTaskNameParam;
	}

	/**
	 * @param dynamicTaskNameParam the dynamicTaskNameParam to set to be used by DYNAMIC tasks
	 *
	 */
	public void setDynamicTaskNameParam(String dynamicTaskNameParam) {
		this.dynamicTaskNameParam = dynamicTaskNameParam;
	}

	/**
	 * @return the caseValueParam
	 */
	public String getCaseValueParam() {
		return caseValueParam;
	}

	@Deprecated
	public String getDynamicForkJoinTasksParam() {
		return dynamicForkJoinTasksParam;
	}

	@Deprecated
	public void setDynamicForkJoinTasksParam(String dynamicForkJoinTasksParam) {
		this.dynamicForkJoinTasksParam = dynamicForkJoinTasksParam;
	}

	public String getDynamicForkTasksParam() {
		return dynamicForkTasksParam;
	}

	public void setDynamicForkTasksParam(String dynamicForkTasksParam) {
		this.dynamicForkTasksParam = dynamicForkTasksParam;
	}

	public String getDynamicForkTasksInputParamName() {
		return dynamicForkTasksInputParamName;
	}

	public void setDynamicForkTasksInputParamName(String dynamicForkTasksInputParamName) {
		this.dynamicForkTasksInputParamName = dynamicForkTasksInputParamName;
	}

	/**
	 * @param caseValueParam the caseValueParam to set
	 */
	public void setCaseValueParam(String caseValueParam) {
		this.caseValueParam = caseValueParam;
	}

	/**
	 *
	 * @return A javascript expression for decision cases.  The result should be a scalar value that is used to decide the case branches.
	 * @see #getDecisionCases()
	 */
	public String getCaseExpression() {
		return caseExpression;
	}

	/**
	 *
	 * @param caseExpression A javascript expression for decision cases.  The result should be a scalar value that is used to decide the case branches.
	 */
	public void setCaseExpression(String caseExpression) {
		this.caseExpression = caseExpression;
	}


	public String getScriptExpression() {
		return scriptExpression;
	}

	public void setScriptExpression(String expression) {
		this.scriptExpression = expression;
	}


	/**
	 * @return the subWorkflow
	 */
	public SubWorkflowParams getSubWorkflowParam() {
		return subWorkflowParam;
	}

	/**
	 * @param subWorkflow the subWorkflowParam to set
	 */
	public void setSubWorkflowParam(SubWorkflowParams subWorkflow) {
		this.subWorkflowParam = subWorkflow;
	}

	/**
	 * @return the joinOn
	 */
	public List<String> getJoinOn() {
		return joinOn;
	}

	/**
	 * @param joinOn the joinOn to set
	 */
	public void setJoinOn(List<String> joinOn) {
		this.joinOn = joinOn;
	}

	/**
	 * @return the loopCondition
	 */
	public String getLoopCondition() {
		return loopCondition;
	}

	/**
	 * @param loopCondition the expression to set
	 */
	public void setLoopCondition(String loopCondition) {
		this.loopCondition = loopCondition;
	}

	/**
	 * @return the loopOver
	 */
	public List<WorkflowTask> getLoopOver() {
		return loopOver;
	}

	/**
	 * @param loopOver the loopOver to set
	 */
	public void setLoopOver(List<WorkflowTask> loopOver) {
		this.loopOver = loopOver;
	}

	/**
	 *
	 * @return Sink value for the EVENT type of task
	 */
	public String getSink() {
		return sink;
	}

	/**
	 *
	 * @param sink Name of the sink
	 */
	public void setSink(String sink) {
		this.sink = sink;
	}

	/**
	 *
	 * @return whether wait for an external event to complete the task, for EVENT and HTTP tasks
	 */
	public Boolean isAsyncComplete() {
		return asyncComplete;
	}

	public void setAsyncComplete(Boolean asyncComplete) {
		this.asyncComplete = asyncComplete;
	}

	/**
	 *
	 * @return If the task is optional.  When set to true, the workflow execution continues even when the task is in failed status.
	 */
	public boolean isOptional() {
		return optional;
	}

	/**
	 *
	 * @return Task definition associated to the Workflow Task
	 */
	public TaskDef getTaskDefinition() {
		return taskDefinition;
	}

	/**
	 * @param taskDefinition Task definition
	 */
	public void setTaskDefinition(TaskDef taskDefinition) {
		this.taskDefinition = taskDefinition;
	}

	/**
	 *
	 * @param optional when set to true, the task is marked as optional
	 */
	public void setOptional(boolean optional) {
		this.optional = optional;
	}

	public Boolean getRateLimited() {
		return rateLimited;
	}

	public void setRateLimited(Boolean rateLimited) {
		this.rateLimited = rateLimited;
	}

	public Boolean isRateLimited() {
		return rateLimited != null && rateLimited;
	}

	public List<String> getDefaultExclusiveJoinTask() {
		return defaultExclusiveJoinTask;
	}

	public void setDefaultExclusiveJoinTask(List<String> defaultExclusiveJoinTask) {
		this.defaultExclusiveJoinTask = defaultExclusiveJoinTask;
	}

	@Override
	public String toString() {
		return name + "/" + taskReferenceName;
	}

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        WorkflowTask that = (WorkflowTask) o;
        return getStartDelay() == that.getStartDelay() &&
                isOptional() == that.isOptional() &&
                Objects.equals(getName(), that.getName()) &&
                Objects.equals(getTaskReferenceName(), that.getTaskReferenceName()) &&
                Objects.equals(getDescription(), that.getDescription()) &&
                Objects.equals(getInputParameters(), that.getInputParameters()) &&
                Objects.equals(getType(), that.getType()) &&
                Objects.equals(getDynamicTaskNameParam(), that.getDynamicTaskNameParam()) &&
                Objects.equals(getCaseValueParam(), that.getCaseValueParam()) &&
                Objects.equals(getCaseExpression(), that.getCaseExpression()) &&
                Objects.equals(getDecisionCases(), that.getDecisionCases()) &&
                Objects.equals(getDynamicForkJoinTasksParam(), that.getDynamicForkJoinTasksParam()) &&
                Objects.equals(getDynamicForkTasksParam(), that.getDynamicForkTasksParam()) &&
                Objects.equals(getDynamicForkTasksInputParamName(), that.getDynamicForkTasksInputParamName()) &&
                Objects.equals(getDefaultCase(), that.getDefaultCase()) &&
                Objects.equals(getForkTasks(), that.getForkTasks()) &&
                Objects.equals(getSubWorkflowParam(), that.getSubWorkflowParam()) &&
                Objects.equals(getJoinOn(), that.getJoinOn()) &&
                Objects.equals(getSink(), that.getSink()) &&
				Objects.equals(isAsyncComplete(), that.isAsyncComplete()) &&
                Objects.equals(getDefaultExclusiveJoinTask(), that.getDefaultExclusiveJoinTask());
    }

    @Override
    public int hashCode() {

        return Objects.hash(
                getName(),
                getTaskReferenceName(),
                getDescription(),
                getInputParameters(),
                getType(),
                getDynamicTaskNameParam(),
                getCaseValueParam(),
                getCaseExpression(),
                getDecisionCases(),
                getDynamicForkJoinTasksParam(),
                getDynamicForkTasksParam(),
                getDynamicForkTasksInputParamName(),
                getDefaultCase(),
                getForkTasks(),
                getStartDelay(),
                getSubWorkflowParam(),
                getJoinOn(),
                getSink(),
                isAsyncComplete(),
                isOptional(),
                getDefaultExclusiveJoinTask()
        );
    }
}
