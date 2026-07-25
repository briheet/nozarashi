CREATE TABLE conversations (
  id UUID PRIMARY KEY,
  title VARCHAR(160) NOT NULL,
  status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'cancelled')),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  cancelled_at TIMESTAMPTZ NULL
);

CREATE INDEX conversations_updated_at_idx ON conversations (updated_at DESC);

CREATE TABLE chat_messages (
  id UUID PRIMARY KEY,
  conversation_id UUID NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
  role VARCHAR(20) NOT NULL CHECK (role IN ('system', 'user', 'assistant')),
  content TEXT NOT NULL CHECK (length(content) > 0),
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX chat_messages_conversation_idx ON chat_messages (conversation_id, created_at);

CREATE TABLE inference_logs (
  id UUID PRIMARY KEY,
  conversation_id UUID NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
  request_message_id UUID NOT NULL REFERENCES chat_messages (id) ON DELETE CASCADE,
  response_message_id UUID NULL,
  provider VARCHAR(40) NOT NULL,
  model VARCHAR(100) NOT NULL,
  status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'error')),
  started_at TIMESTAMPTZ NOT NULL,
  completed_at TIMESTAMPTZ NOT NULL,
  latency_ms BIGINT NOT NULL CHECK (latency_ms >= 0),
  input_tokens BIGINT NOT NULL DEFAULT 0 CHECK (input_tokens >= 0),
  output_tokens BIGINT NOT NULL DEFAULT 0 CHECK (output_tokens >= 0),
  total_tokens BIGINT NOT NULL DEFAULT 0 CHECK (total_tokens >= 0),
  input_preview VARCHAR(1000) NOT NULL DEFAULT '',
  output_preview VARCHAR(1000) NOT NULL DEFAULT '',
  error_message VARCHAR(2000) NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX inference_logs_started_at_idx ON inference_logs (started_at DESC);
CREATE INDEX inference_logs_model_status_idx ON inference_logs (model, status, started_at DESC);
CREATE INDEX inference_logs_conversation_started_at_id_idx
  ON inference_logs (conversation_id, started_at, id);
