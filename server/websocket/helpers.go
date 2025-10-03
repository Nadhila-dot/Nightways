package websocket

import "time"

// Fix these soon
// To be proper

func Push(message string, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "type": "push",
        "data": map[string]interface{}{
            "type":    "push",
            "message": message,
            "extra":   extra,
        },
    }
}

func Start(message string, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "data": map[string]interface{}{
            "type":    "start",
            "message": message,
            "extra":   extra,
        },
    }
}

func Retry(message string, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "data": map[string]interface{}{
            "type":    "retry",
            "message": message,
            "extra":   extra,
        },
    }
}

func End(message string, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "data": map[string]interface{}{
            "type":    "end",
            "message": message,
            "time":   time.Now().Format(time.RFC3339),
            "extra":   extra,
        },
    }
}

func Complete(message string,extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "data": map[string]interface{}{
            "type":  "complete",
            "message": message,
            "extra": extra,
        },
    }
}

func Progression(message string, current, total int, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "data": map[string]interface{}{
            "type":    "progress",
            "message": message,
            "progress": map[string]interface{}{
                "current": current,
                "total":   total,
            },
            "extra": extra,
        },
    }
}

func Stage(stage string, step string, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "data": map[string]interface{}{
            "type":  "stage",
            "stage": stage,
            "step":  step,
            "extra": extra,
        },
    }
}

func Review_output(heading string, content string, need bool, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
       
        "data": map[string]interface{}{
            "type":  "review-out",
            "modal": map[string]interface{}{
                "heading": heading,
                // The content should support Markdown
                // I think that's better
                "content":  content,
                "optional": need,
            },
            "extra": extra,
        },
    }
}

func Progress(message string, current, total int, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "type": "progress",
        "data": map[string]interface{}{
            "type":    "progress",
            "message": message,
            "progress": map[string]interface{}{
                "current": current,
                "total":   total,
            },
            "extra": extra,
        },
    }
}

func Completed(message string, result interface{}, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "type": "completed",
        "data": map[string]interface{}{
            "type":    "completed",
            "message": message,
            "result":  result,
            "extra":   extra,
        },
    }
}

func Error(message string, err string, extra map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "type": "error",
        "data": map[string]interface{}{
            "type":    "error",
            "message": message,
            "error":   err,
            "extra":   extra,
        },
    }
}

func Info(infoType, message string) map[string]interface{} {
    return map[string]interface{}{
        "data": map[string]interface{}{
            "type":    infoType,
            "message": message,
        },
    }
}