CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    is_archived BOOLEAN NOT NULL DEFAULT false,
    parent_id UUID,
    content TEXT,
    user_id UUID NOT NULL REFERENCES users(id),
    document_path VARCHAR(255) NOT NULL,
    workspace_id UUID NOT NULL REFERENCES workspaces(id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_documents_user_id FOREIGN KEY (user_id) REFERENCES users(id)
    CONSTRAINT fk_documents_parent_id FOREIGN KEY (parent_id) REFERENCES documents(id)
    CONSTRAINT fk_documents_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
);

CREATE INDEX IF NOT EXISTS idx_documents_document_path ON documents(document_path);
CREATE INDEX IF NOT EXISTS idx_documents_parent_id ON documents(parent_id);
CREATE INDEX IF NOT EXISTS idx_documents_workspace_id ON documents(workspace_id);