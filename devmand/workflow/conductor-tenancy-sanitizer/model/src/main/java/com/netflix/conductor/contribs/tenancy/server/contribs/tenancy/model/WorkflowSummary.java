package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model;

import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model.Workflow.WorkflowStatus;
import java.util.Objects;
import java.util.TimeZone;

/**
 * Captures workflow summary info to be indexed in Elastic Search.
 *
 * @author Viren
 */
public class WorkflowSummary {

	/**
	 * The time should be stored as GMT
	 */
	private static final TimeZone gmt = TimeZone.getTimeZone("GMT");


	private String workflowType;


	private int version;


	private String workflowId;


	private String correlationId;


	private String startTime;


	private String updateTime;


	private String endTime;


	private WorkflowStatus status;


	private String input;


	private String output;


	private String reasonForIncompletion;


	private long executionTime;


	private String event;


	private String failedReferenceTaskNames = "";


	private String externalInputPayloadStoragePath;


	private String externalOutputPayloadStoragePath;


	private int priority;
	private long inputSize;
	private long outputSize;


	/**
	 * @return the workflowType
	 */
	public String getWorkflowType() {
		return workflowType;
	}

	/**
	 * @return the version
	 */
	public int getVersion() {
		return version;
	}

	/**
	 * @return the workflowId
	 */
	public String getWorkflowId() {
		return workflowId;
	}

	/**
	 * @return the correlationId
	 */
	public String getCorrelationId() {
		return correlationId;
	}

	/**
	 * @return the startTime
	 */
	public String getStartTime() {
		return startTime;
	}

	/**
	 * @return the endTime
	 */
	public String getEndTime() {
		return endTime;
	}

	/**
	 * @return the status
	 */
	public WorkflowStatus getStatus() {
		return status;
	}

	/**
	 * @return the input
	 */
	public String getInput() {
		return input;
	}


    public long getInputSize() {
        return inputSize;
    }

	public void setInputSize(long inputSize) {
		this.inputSize = inputSize;
	}

	/**
	 *
	 * @return the output
	 */
	public String getOutput() {
		return output;
	}

    public long getOutputSize() {
        return outputSize;
    }

	public void setOutputSize(long outputSize) {
		this.outputSize = outputSize;
	}

	/**
	 * @return the reasonForIncompletion
	 */
	public String getReasonForIncompletion() {
		return reasonForIncompletion;
	}

	/**
	 *
	 * @return the executionTime
	 */
	public long getExecutionTime(){
		return executionTime;
	}

	/**
	 * @return the updateTime
	 */
	public String getUpdateTime() {
		return updateTime;
	}

	/**
	 *
	 * @return The event
	 */
	public String getEvent() {
		return event;
	}

	/**
	 *
	 * @param event The event
	 */
	public void setEvent(String event) {
		this.event = event;
	}

	public String getFailedReferenceTaskNames() {
		return failedReferenceTaskNames;
	}

	public void setFailedReferenceTaskNames(String failedReferenceTaskNames) {
		this.failedReferenceTaskNames = failedReferenceTaskNames;
	}

	public void setWorkflowType(String workflowType) {
		this.workflowType = workflowType;
	}

	public void setVersion(int version) {
		this.version = version;
	}

	public void setWorkflowId(String workflowId) {
		this.workflowId = workflowId;
	}

	public void setCorrelationId(String correlationId) {
		this.correlationId = correlationId;
	}

	public void setStartTime(String startTime) {
		this.startTime = startTime;
	}

	public void setUpdateTime(String updateTime) {
		this.updateTime = updateTime;
	}

	public void setEndTime(String endTime) {
		this.endTime = endTime;
	}

	public void setStatus(WorkflowStatus status) {
		this.status = status;
	}

	public void setInput(String input) {
		this.input = input;
	}

	public void setOutput(String output) {
		this.output = output;
	}

	public void setReasonForIncompletion(String reasonForIncompletion) {
		this.reasonForIncompletion = reasonForIncompletion;
	}

	public void setExecutionTime(long executionTime) {
		this.executionTime = executionTime;
	}

	/**
	 * @return the external storage path of the workflow input payload
	 */
	public String getExternalInputPayloadStoragePath() {
		return externalInputPayloadStoragePath;
	}

	/**
	 * @param externalInputPayloadStoragePath the external storage path where the workflow input payload is stored
	 */
	public void setExternalInputPayloadStoragePath(String externalInputPayloadStoragePath) {
		this.externalInputPayloadStoragePath = externalInputPayloadStoragePath;
	}

	/**
	 * @return the external storage path of the workflow output payload
	 */
	public String getExternalOutputPayloadStoragePath() {
		return externalOutputPayloadStoragePath;
	}

	/**
	 * @param externalOutputPayloadStoragePath the external storage path where the workflow output payload is stored
	 */
	public void setExternalOutputPayloadStoragePath(String externalOutputPayloadStoragePath) {
		this.externalOutputPayloadStoragePath = externalOutputPayloadStoragePath;
	}

	/**
	 * @return the priority to define on tasks
	 */
	public int getPriority() {
		return priority;
	}

	/**
	 * @param priority priority of tasks (between 0 and 99)
	 */
	public void setPriority(int priority) {
		this.priority = priority;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) {
			return true;
		}
		if (o == null || getClass() != o.getClass()) {
			return false;
		}
		WorkflowSummary that = (WorkflowSummary) o;
		return getVersion() == that.getVersion() &&
			getExecutionTime() == that.getExecutionTime() &&
			getPriority() == that.getPriority() &&
			getWorkflowType().equals(that.getWorkflowType()) &&
			getWorkflowId().equals(that.getWorkflowId()) &&
			Objects.equals(getCorrelationId(), that.getCorrelationId()) &&
			getStartTime().equals(that.getStartTime()) &&
			getUpdateTime().equals(that.getUpdateTime()) &&
			getEndTime().equals(that.getEndTime()) &&
			getStatus() == that.getStatus() &&
			Objects.equals(getReasonForIncompletion(), that.getReasonForIncompletion()) &&
			Objects.equals(getEvent(), that.getEvent());
	}

	@Override
	public int hashCode() {
		return Objects
			.hash(getWorkflowType(), getVersion(), getWorkflowId(), getCorrelationId(), getStartTime(), getUpdateTime(),
				getEndTime(), getStatus(), getReasonForIncompletion(), getExecutionTime(), getEvent(), getPriority());
	}
}
