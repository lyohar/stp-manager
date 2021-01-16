CREATE TABLE public.exports
(
    table_schema          text NOT NULL,
    table_name            text NOT NULL,
    columns_list          text,
    keys_list              text,
    order_column_name     text,
    order_column_value    text,
    timestamp_column_name text,
    operation_column_name text,
    topic_name            text,
    executing             int4 NOT NULL,
    CONSTRAINT exports_pkey PRIMARY KEY (table_schema, table_name)
);