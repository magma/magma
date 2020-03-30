package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model;

import java.util.Map;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;

/**
 * @author Viren
 *
 */

public class SubWorkflowParams {


    @NotNull(message = "SubWorkflowParams name cannot be null")
    @NotEmpty(message = "SubWorkflowParams name cannot be empty")
    private String name;


    private Integer version;


    private Map<String, String> taskToDomain;

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
     * @return the version
     */
    public Integer getVersion() {
        return version;
    }

    /**
     * @param version the version to set
     */
    public void setVersion(Integer version) {
        this.version = version;
    }

    /**
     * @return the taskToDomain
     */
    public Map<String, String> getTaskToDomain() {
        return taskToDomain;
    }
    /**
     * @param taskToDomain the taskToDomain to set
     */
    public void setTaskToDomain(Map<String, String> taskToDomain) {
        this.taskToDomain = taskToDomain;
    }
}
