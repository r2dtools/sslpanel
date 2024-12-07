export enum FetchStatus {
    Idle = 'idle',
    Pending = 'pending',
    Failed = 'failed',
    Succeeded = 'succeeded',
};

export interface RouteItem {
    path?: string;
    index?: boolean;
    public?: boolean;
    title: string;
    name?: string;
    component: React.ReactNode | null;
};
