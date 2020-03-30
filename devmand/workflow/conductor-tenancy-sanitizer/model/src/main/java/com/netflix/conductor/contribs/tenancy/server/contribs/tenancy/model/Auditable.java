package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model;

/**
 * @author Viren
 *
 */
public abstract class Auditable {

	private String ownerApp;

	private Long createTime;

	private Long updateTime;

	private String createdBy;

	private String updatedBy;

	/**
	 * @return the ownerApp
	 */
	public String getOwnerApp() {
		return ownerApp;
	}

	/**
	 * @param ownerApp the ownerApp to set
	 */
	public void setOwnerApp(String ownerApp) {
		this.ownerApp = ownerApp;
	}

	/**
	 * @return the createTime
	 */
	public Long getCreateTime() {
		return createTime;
	}

	/**
	 * @param createTime the createTime to set
	 */
	public void setCreateTime(Long createTime) {
		this.createTime = createTime;
	}

	/**
	 * @return the updateTime
	 */
	public Long getUpdateTime() {
		return updateTime;
	}

	/**
	 * @param updateTime the updateTime to set
	 */
	public void setUpdateTime(Long updateTime) {
		this.updateTime = updateTime;
	}

	/**
	 * @return the createdBy
	 */
	public String getCreatedBy() {
		return createdBy;
	}

	/**
	 * @param createdBy the createdBy to set
	 */
	public void setCreatedBy(String createdBy) {
		this.createdBy = createdBy;
	}

	/**
	 * @return the updatedBy
	 */
	public String getUpdatedBy() {
		return updatedBy;
	}

	/**
	 * @param updatedBy the updatedBy to set
	 */
	public void setUpdatedBy(String updatedBy) {
		this.updatedBy = updatedBy;
	}
}
