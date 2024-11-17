import { useEffect } from 'react';
import Breadcrumb from '../../components/Breadcrumb';
import ServerItem from '../../features/server/components/ServerItem';
import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { fetchServers, selectServers, selectServersFetchStatus } from '../../features/server/serverSlice';
import useAuthToken from '../../features/auth/hooks';
import { FetchStatus } from '../../app/types';
import Loader from '../../components/Loader/Loader';
import { Button } from 'flowbite-react';

const ServerList = () => {
    const [authToken] = useAuthToken();
    const dispatch = useAppDispatch();
    const serversSelectStatus = useAppSelector(selectServersFetchStatus);
    const servers = useAppSelector(selectServers);

    useEffect(() => {
        if (authToken && serversSelectStatus !== FetchStatus.Pending) {
            dispatch(fetchServers(authToken));
        }
    }, [authToken]);

    return (
        serversSelectStatus !== FetchStatus.Pending ?
            <>
                <Breadcrumb pageName='Servers'>
                    <Button color='blue'>Add</Button>
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
                                    <ServerItem key={server.id} server={server} />
                                ))
                            }
                        </div>
                    </div>
                </div>
            </> : <Loader />
    );
}

export default ServerList;
