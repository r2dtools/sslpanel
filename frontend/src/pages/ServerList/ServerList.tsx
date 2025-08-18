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
import { Banner, BannerCollapseButton, Button } from 'flowbite-react';
import { Server, ServerSavePayload } from '../../features/server/types';
import ServerEditDrawer from '../../features/server/components/ServerEditDrawer';
import empty from '../../images/empty.png';
import { HiX } from "react-icons/hi";
import { Link } from 'react-router-dom';
import { HiMiniExclamationCircle } from 'react-icons/hi2';
import Error404 from '../Error404';

const ServerList = () => {
    const [serverFormOpen, setServerFormOpen] = useState(false);
    const [editedServer, setEditedServer] = useState<Server | null>(null);
    const [authToken] = useAuthToken();
    const dispatch = useAppDispatch();
    const serversSelectStatus = useAppSelector(selectServersFetchStatus);
    const serverSaveStatus = useAppSelector(selectServerSaveStatus);
    const servers = useAppSelector(selectServers);

    if (!authToken) {
        return <Error404 />
    }

    useEffect(() => {
        if (serversSelectStatus !== FetchStatus.Pending) {
            dispatch(fetchServers(authToken));
        }
    }, [authToken]);

    // close drawer on success
    useEffect(() => {
        if (serverFormOpen && serverSaveStatus === FetchStatus.Succeeded) {
            handleServerFormClose();
        }
    }, [serverSaveStatus]);

    const handleServerFormClose = (): void => {
        setServerFormOpen(false);
        setEditedServer(null);
    };

    const handleServerFormOpen = (): void => {
        setServerFormOpen(true);
    };

    const handleSubmit = async (server: ServerSavePayload) => {
        server.id ? dispatch(editServer(server)) : await dispatch(addServer(server));
    };

    const handleDeleteServer = (id: number) => {
        return dispatch(deleteServer({ id, token: authToken }));
    };

    const handleEditServer = (server: Server) => {
        setEditedServer(server);
        handleServerFormOpen();
    };

    return (
        serversSelectStatus !== FetchStatus.Pending ?
            <>
                <Breadcrumb pageName='Servers' hideNavigation>
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
                        <div className="mx-auto max-w-[80px]">
                            {!servers.length && <img src={empty} />}
                        </div>
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
                    {!servers.length &&
                        <Banner>
                            <div className="flex w-full justify-between border-b border-gray-200 bg-gray-50 p-4 dark:border-gray-600 dark:bg-gray-700">
                                <div className="mx-auto flex items-center">
                                    <p className="flex items-center font-normal text-gray-500 dark:text-gray-400">
                                        <HiMiniExclamationCircle className="mr-2 h-5 w-5" />
                                        <span className="[&_p]:inline">
                                            Hi! Looks like you haven't registered any servers yet. Follow the&nbsp;
                                            <Link to='/documentation/install-sslbot' className="text-blue-600 inline font-medium underline decoration-solid underline-offset-2 hover:no-underlin">instructions</Link>
                                            &nbsp;to install SSLBot and register a server
                                        </span>
                                    </p>
                                </div>
                                <BannerCollapseButton color="gray" className="border-0 bg-transparent text-gray-500 dark:text-gray-400">
                                    <HiX className="h-4 w-4" />
                                </BannerCollapseButton>
                            </div>
                        </Banner>
                    }
                </div>
                <ServerEditDrawer
                    open={serverFormOpen}
                    authToken={authToken}
                    loading={serverSaveStatus === FetchStatus.Pending}
                    onSubmit={handleSubmit}
                    onClose={handleServerFormClose}
                    server={editedServer}
                />
            </> : <Loader />
    );
}

export default ServerList;
