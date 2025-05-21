package kbapi

type FleetListAgentsResponseBodyAgents struct {
	ID                   string                 `json:"id"`
	Version              string                 `json:"version"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type FleetListAgentsResponseBody struct {
	Items []struct {
		AccessApiKey          *string                            `json:"access_api_key,omitempty"`
		AccessApiKeyId        *string                            `json:"access_api_key_id,omitempty"`
		Active                bool                               `json:"active"`
		Agent                 *FleetListAgentsResponseBodyAgents `json:"agent,omitempty"`
		AuditUnenrolledReason *string                            `json:"audit_unenrolled_reason,omitempty"`
		Components            *[]struct {
			Id      string `json:"id"`
			Message string `json:"message"`
			Status  string `json:"status"`
			Type    string `json:"type"`
			Units   *[]struct {
				Id      string                  `json:"id"`
				Message string                  `json:"message"`
				Payload *map[string]interface{} `json:"payload,omitempty"`
				Status  string                  `json:"status"`
				Type    string                  `json:"type"`
			} `json:"units,omitempty"`
		} `json:"components,omitempty"`
		DefaultApiKey        *string `json:"default_api_key,omitempty"`
		DefaultApiKeyHistory *[]struct {
			Id        string `json:"id"`
			RetiredAt string `json:"retired_at"`
		} `json:"default_api_key_history,omitempty"`
		DefaultApiKeyId    *string                `json:"default_api_key_id,omitempty"`
		EnrolledAt         string                 `json:"enrolled_at"`
		Id                 string                 `json:"id"`
		LastCheckin        *string                `json:"last_checkin,omitempty"`
		LastCheckinMessage *string                `json:"last_checkin_message,omitempty"`
		LastCheckinStatus  *string                `json:"last_checkin_status,omitempty"`
		LocalMetadata      map[string]interface{} `json:"local_metadata"`
		Metrics            *struct {
			CpuAvg            *float32 `json:"cpu_avg,omitempty"`
			MemorySizeByteAvg *float32 `json:"memory_size_byte_avg,omitempty"`
		} `json:"metrics,omitempty"`
		Namespaces *[]string `json:"namespaces,omitempty"`
		Outputs    *map[string]struct {
			ApiKeyId          string `json:"api_key_id"`
			ToRetireApiKeyIds *[]struct {
				Id        string `json:"id"`
				RetiredAt string `json:"retired_at"`
			} `json:"to_retire_api_key_ids,omitempty"`
			Type string `json:"type"`
		} `json:"outputs,omitempty"`
		Packages              []string       `json:"packages"`
		PolicyId              *string        `json:"policy_id,omitempty"`
		PolicyRevision        *float32       `json:"policy_revision"`
		Sort                  *[]interface{} `json:"sort,omitempty"`
		Status                *string        `json:"status,omitempty"`
		Tags                  *[]string      `json:"tags,omitempty"`
		Type                  string         `json:"type"`
		UnenrolledAt          *string        `json:"unenrolled_at,omitempty"`
		UnenrollmentStartedAt *string        `json:"unenrollment_started_at,omitempty"`
		UnhealthyReason       *[]string      `json:"unhealthy_reason"`
		UpgradeAttempts       *[]string      `json:"upgrade_attempts"`
		UpgradeDetails        *struct {
			ActionId string `json:"action_id"`
			Metadata *struct {
				DownloadPercent *float32 `json:"download_percent,omitempty"`
				DownloadRate    *float32 `json:"download_rate,omitempty"`
				ErrorMsg        *string  `json:"error_msg,omitempty"`
				FailedState     *string  `json:"failed_state,omitempty"`
				RetryErrorMsg   *string  `json:"retry_error_msg,omitempty"`
				RetryUntil      *string  `json:"retry_until,omitempty"`
				ScheduledAt     *string  `json:"scheduled_at,omitempty"`
			} `json:"metadata,omitempty"`
			State         string `json:"state"`
			TargetVersion string `json:"target_version"`
		} `json:"upgrade_details"`
		UpgradeStartedAt     *string                 `json:"upgrade_started_at"`
		UpgradedAt           *string                 `json:"upgraded_at"`
		UserProvidedMetadata *map[string]interface{} `json:"user_provided_metadata,omitempty"`
	} `json:"items"`
	NextSearchAfter *string             `json:"nextSearchAfter,omitempty"`
	Page            float32             `json:"page"`
	PerPage         float32             `json:"perPage"`
	Pit             *string             `json:"pit,omitempty"`
	StatusSummary   *map[string]float32 `json:"statusSummary,omitempty"`
	Total           float32             `json:"total"`
}

type FleetGetAgentResponseBody struct {
	Item struct {
		AccessApiKey          *string                            `json:"access_api_key,omitempty"`
		AccessApiKeyId        *string                            `json:"access_api_key_id,omitempty"`
		Active                bool                               `json:"active"`
		Agent                 *FleetListAgentsResponseBodyAgents `json:"agent,omitempty"`
		AuditUnenrolledReason *string                            `json:"audit_unenrolled_reason,omitempty"`
		Components            *[]struct {
			Id      string `json:"id"`
			Message string `json:"message"`
			Status  string `json:"status"`
			Type    string `json:"type"`
			Units   *[]struct {
				Id      string                  `json:"id"`
				Message string                  `json:"message"`
				Payload *map[string]interface{} `json:"payload,omitempty"`
				Status  string                  `json:"status"`
				Type    string                  `json:"type"`
			} `json:"units,omitempty"`
		} `json:"components,omitempty"`
		DefaultApiKey        *string `json:"default_api_key,omitempty"`
		DefaultApiKeyHistory *[]struct {
			Id        string `json:"id"`
			RetiredAt string `json:"retired_at"`
		} `json:"default_api_key_history,omitempty"`
		DefaultApiKeyId    *string                `json:"default_api_key_id,omitempty"`
		EnrolledAt         string                 `json:"enrolled_at"`
		Id                 string                 `json:"id"`
		LastCheckin        *string                `json:"last_checkin,omitempty"`
		LastCheckinMessage *string                `json:"last_checkin_message,omitempty"`
		LastCheckinStatus  *string                `json:"last_checkin_status,omitempty"`
		LocalMetadata      map[string]interface{} `json:"local_metadata"`
		Metrics            *struct {
			CpuAvg            *float32 `json:"cpu_avg,omitempty"`
			MemorySizeByteAvg *float32 `json:"memory_size_byte_avg,omitempty"`
		} `json:"metrics,omitempty"`
		Namespaces *[]string `json:"namespaces,omitempty"`
		Outputs    *map[string]struct {
			ApiKeyId          string `json:"api_key_id"`
			ToRetireApiKeyIds *[]struct {
				Id        string `json:"id"`
				RetiredAt string `json:"retired_at"`
			} `json:"to_retire_api_key_ids,omitempty"`
			Type string `json:"type"`
		} `json:"outputs,omitempty"`
		Packages              []string       `json:"packages"`
		PolicyId              *string        `json:"policy_id,omitempty"`
		PolicyRevision        *float32       `json:"policy_revision"`
		Sort                  *[]interface{} `json:"sort,omitempty"`
		Status                *string        `json:"status,omitempty"`
		Tags                  *[]string      `json:"tags,omitempty"`
		Type                  string         `json:"type"`
		UnenrolledAt          *string        `json:"unenrolled_at,omitempty"`
		UnenrollmentStartedAt *string        `json:"unenrollment_started_at,omitempty"`
		UnhealthyReason       *[]string      `json:"unhealthy_reason"`
		UpgradeAttempts       *[]string      `json:"upgrade_attempts"`
		UpgradeDetails        *struct {
			ActionId string `json:"action_id"`
			Metadata *struct {
				DownloadPercent *float32 `json:"download_percent,omitempty"`
				DownloadRate    *float32 `json:"download_rate,omitempty"`
				ErrorMsg        *string  `json:"error_msg,omitempty"`
				FailedState     *string  `json:"failed_state,omitempty"`
				RetryErrorMsg   *string  `json:"retry_error_msg,omitempty"`
				RetryUntil      *string  `json:"retry_until,omitempty"`
				ScheduledAt     *string  `json:"scheduled_at,omitempty"`
			} `json:"metadata,omitempty"`
			State         string `json:"state"`
			TargetVersion string `json:"target_version"`
		} `json:"upgrade_details"`
		UpgradeStartedAt     *string                 `json:"upgrade_started_at"`
		UpgradedAt           *string                 `json:"upgraded_at"`
		UserProvidedMetadata *map[string]interface{} `json:"user_provided_metadata,omitempty"`
	} `json:"item"`
}

type FleetUpdateAgentResponseBody struct {
	Item struct {
		AccessApiKey          *string                            `json:"access_api_key,omitempty"`
		AccessApiKeyId        *string                            `json:"access_api_key_id,omitempty"`
		Active                bool                               `json:"active"`
		Agent                 *FleetListAgentsResponseBodyAgents `json:"agent,omitempty"`
		AuditUnenrolledReason *string                            `json:"audit_unenrolled_reason,omitempty"`
		Components            *[]struct {
			Id      string `json:"id"`
			Message string `json:"message"`
			Status  string `json:"status"`
			Type    string `json:"type"`
			Units   *[]struct {
				Id      string                  `json:"id"`
				Message string                  `json:"message"`
				Payload *map[string]interface{} `json:"payload,omitempty"`
				Status  string                  `json:"status"`
				Type    string                  `json:"type"`
			} `json:"units,omitempty"`
		} `json:"components,omitempty"`
		DefaultApiKey        *string `json:"default_api_key,omitempty"`
		DefaultApiKeyHistory *[]struct {
			Id        string `json:"id"`
			RetiredAt string `json:"retired_at"`
		} `json:"default_api_key_history,omitempty"`
		DefaultApiKeyId    *string                `json:"default_api_key_id,omitempty"`
		EnrolledAt         string                 `json:"enrolled_at"`
		Id                 string                 `json:"id"`
		LastCheckin        *string                `json:"last_checkin,omitempty"`
		LastCheckinMessage *string                `json:"last_checkin_message,omitempty"`
		LastCheckinStatus  *string                `json:"last_checkin_status,omitempty"`
		LocalMetadata      map[string]interface{} `json:"local_metadata"`
		Metrics            *struct {
			CpuAvg            *float32 `json:"cpu_avg,omitempty"`
			MemorySizeByteAvg *float32 `json:"memory_size_byte_avg,omitempty"`
		} `json:"metrics,omitempty"`
		Namespaces *[]string `json:"namespaces,omitempty"`
		Outputs    *map[string]struct {
			ApiKeyId          string `json:"api_key_id"`
			ToRetireApiKeyIds *[]struct {
				Id        string `json:"id"`
				RetiredAt string `json:"retired_at"`
			} `json:"to_retire_api_key_ids,omitempty"`
			Type string `json:"type"`
		} `json:"outputs,omitempty"`
		Packages              []string       `json:"packages"`
		PolicyId              *string        `json:"policy_id,omitempty"`
		PolicyRevision        *float32       `json:"policy_revision"`
		Sort                  *[]interface{} `json:"sort,omitempty"`
		Status                *string        `json:"status,omitempty"`
		Tags                  *[]string      `json:"tags,omitempty"`
		Type                  string         `json:"type"`
		UnenrolledAt          *string        `json:"unenrolled_at,omitempty"`
		UnenrollmentStartedAt *string        `json:"unenrollment_started_at,omitempty"`
		UnhealthyReason       *[]string      `json:"unhealthy_reason"`
		UpgradeAttempts       *[]string      `json:"upgrade_attempts"`
		UpgradeDetails        *struct {
			ActionId string `json:"action_id"`
			Metadata *struct {
				DownloadPercent *float32 `json:"download_percent,omitempty"`
				DownloadRate    *float32 `json:"download_rate,omitempty"`
				ErrorMsg        *string  `json:"error_msg,omitempty"`
				FailedState     *string  `json:"failed_state,omitempty"`
				RetryErrorMsg   *string  `json:"retry_error_msg,omitempty"`
				RetryUntil      *string  `json:"retry_until,omitempty"`
				ScheduledAt     *string  `json:"scheduled_at,omitempty"`
			} `json:"metadata,omitempty"`
			State         string `json:"state"`
			TargetVersion string `json:"target_version"`
		} `json:"upgrade_details"`
		UpgradeStartedAt     *string                 `json:"upgrade_started_at"`
		UpgradedAt           *string                 `json:"upgraded_at"`
		UserProvidedMetadata *map[string]interface{} `json:"user_provided_metadata,omitempty"`
	} `json:"item"`
}
