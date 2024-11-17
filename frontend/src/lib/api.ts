import axios, { AxiosError, AxiosRequestConfig } from 'axios';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_HOST,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const getErrorMessage = (error: unknown) => {
    if (error instanceof AxiosError && error.response) {
        const data = error.response.data as { code: number, message: string }

        return data.message;
    }

    if (error instanceof Error) {
        return error.message;
    }

    return String(error)
}

export const configWithAuth = (token: string, config?: AxiosRequestConfig): AxiosRequestConfig => {
    config = config || {};
    const headers = config.headers || {};

    return {
        ...config,
        headers: {
            ...headers,
            Authorization: `Bearer ${token}`,
        },
    };
};


export default api;
