import { useEffect, useState } from 'react';
import Breadcrumb from '../../components/Breadcrumb';
import ServerItem from '../../features/server/components/ServerItem';
import { useAppDispatch, useAppSelector } from '../../app/hooks';
import {
    addServer,
    deleteServer,
    editServer,
    fetchServers,
    selectServers,
    selectServerSaveStatus,
    selectServersFetchStatus,
} from '../../features/server/serversSlice';
import useAuthToken from '../../features/auth/hooks';
import { FetchStatus } from '../../app/types';
import Loader from '../../components/Loader/Loader';
import { Button } from 'flowbite-react';
import { Server, ServerSavePayload } from '../../features/server/types';
import ServerEditDrawer from '../../features/server/components/ServerEditDrawer';

const ServerList = () => {
    const [serverFormOpen, setServerFormOpen] = useState(false);
    const [editedServer, setEditedServer] = useState<Server | null>(null);
    const [authToken] = useAuthToken();
    const dispatch = useAppDispatch();
    const serversSelectStatus = useAppSelector(selectServersFetchStatus);
    const serverSaveStatus = useAppSelector(selectServerSaveStatus);
    const servers = useAppSelector(selectServers);

    useEffect(() => {
        if (authToken && serversSelectStatus !== FetchStatus.Pending) {
            dispatch(fetchServers(authToken));
        }
    }, [authToken]);

    const handleServerFormClose = (): void => {
        setServerFormOpen(false);
        setEditedServer(null);
    };

    const handleServerFormOpen = (): void => {
        setServerFormOpen(true);
    };

    const handleSubmit = async (server: ServerSavePayload) => {
        server.id ? dispatch(editServer(server)) : dispatch(addServer(server));
    };

    const handleDeleteServer = (id: number) => {
        return dispatch(deleteServer({ id, token: authToken || '' }));
    };

    const handleEditServer = (server: Server) => {
        setEditedServer(server);
        handleServerFormOpen();
    };

    return (
        serversSelectStatus !== FetchStatus.Pending ?
            <>
                <Breadcrumb pageName='Servers'>
                    <Button color='blue' onClick={handleServerFormOpen}>Add</Button>
                </Breadcrumb>
                <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                    <div className="border-b border-stroke px-4 py-5 dark:border-strokedark md:px-6 xl:px-7.5">
                        <div className="flex items-center gap-3">
                            <div className="md:w-3/12 w-6/12 lg:w-4/12">
                                <span className="font-medium">Name</span>
                            </div>
                            <div className="hidden md:block md:w-7/12 lg:w-6/12 xl:w-5/12">
                                <span className="font-medium">IP</span>
                            </div>
                            <div className="hidden lg:w-1/12 xl:block">
                                <span className="font-medium">Agent</span>
                            </div>
                            <div className="w-3/12 md:w-1/12">
                                <span className="font-medium">Status</span>
                            </div>
                            <div className="w-3/12 text-center md:w-1/12">
                                <span className="font-medium">Actions</span>
                            </div>
                        </div>
                    </div>

                    <div className="p-4 md:p-6 xl:p-7.5">
                        <div className="flex flex-col gap-7">
                            {
                                servers.map(server => (
                                    <ServerItem
                                        key={server.id}
                                        server={server}
                                        onDelete={handleDeleteServer}
                                        onEdit={handleEditServer}
                                    />
                                ))
                            }
                        </div>
                    </div>
                </div>
                <ServerEditDrawer
                    open={serverFormOpen}
                    authToken={authToken || ''}
                    loading={serverSaveStatus === FetchStatus.Pending}
                    onSubmit={handleSubmit}
                    onClose={handleServerFormClose}
                    server={editedServer}
                />
            </> : <Loader />
    );
}

export default ServerList;
