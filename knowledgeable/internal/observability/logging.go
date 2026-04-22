package observability

func LogAttrs(event, trackingID string, userID *int64, extra ...any) []any {
	isAnonymous := userID == nil

	attrs := []any{
		"event", event,
		"tracking_id", trackingID,
		"is_anonymous", isAnonymous,
	}

	if userID != nil {
		attrs = append(attrs, "user_id", *userID)
	}

	attrs = append(attrs, extra...)
	return attrs
}