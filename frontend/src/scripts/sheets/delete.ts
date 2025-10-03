import http from "@/http";
import axios, { AxiosRequestConfig } from "axios";



export async function deleteJob(jobId: string, session?: string): Promise<{ status: string; error?: string }> {
    const config: AxiosRequestConfig = {
        headers: {},
    };
    // Always use session from localStorage if not provided
    if (!session) {
        session = localStorage.getItem("session") || "";
    }
    if (session) {
        config.headers!["Authorization"] = `Bearer ${session}`;
    }
    try {
        const res = await http.post(`/api/v1/sheets/queue/${jobId}`, config);
        if (res.status < 200 || res.status >= 300) {
            return { status: "error", error: res.data?.error || "Failed to delete job" };
        }
        return { status: "deleted" };
    } catch (e: any) {
        return { status: "error", error: e.response?.data?.error || e.message || "Network error" };
    }
}