import polars as pl
import streamlit as st
from st_supabase_connection import SupabaseConnection

conn = st.connection("supabase", type=SupabaseConnection)

# Use the Supabase client's table API
rows = conn.table("sets").select("*").execute()

data = rows.data

df = pl.from_dicts(data)

st.write(df)
