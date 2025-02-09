import { Alert, Avatar, Badge, Button, Popover, Tooltip } from 'flowbite-react';
import Breadcrumb from '../../components/Breadcrumb';
import { getOsIcon, getOsName } from '../../features/server/utils';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { HiMiniSignal, HiMiniPencil, HiMiniEye, HiMiniClipboard } from 'react-icons/hi2';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import { toast } from 'react-toastify';
import useAuthToken from '../../features/auth/hooks';
import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { useEffect, useState } from 'react';
import {
    editServer,
    fetchServer,
    fetchServerDetails,
    selectServer,
    selectServerDetails,
    selectServerDetailsFetchStatus,
    selectServerFetchStatus,
    selectServerSaveStatus,
} from '../../features/server/serverSlice';
import { FetchStatus } from '../../app/types';
import Loader from '../../components/Loader/Loader';
import moment from 'moment';
import Error404 from '../Error404';
import ServerDomainList from '../../features/server/components/ServerDomainList';
import ServerEditDrawer from '../../features/server/components/ServerEditDrawer';
import { Domain, ServerSavePayload } from '../../features/server/types';
import { domainFetched } from '../../features/domain/domainSlice';
import { encode } from 'js-base64';

const emptyPlaceholder = '----------';

const Server = () => {
    const { guid } = useParams();
    const [authToken] = useAuthToken();
    const dispatch = useAppDispatch();
    const serverSelectStatus = useAppSelector(selectServerFetchStatus);
    const serverDetailsSelectStatus = useAppSelector(selectServerDetailsFetchStatus);
    const server = useAppSelector(selectServer);
    const serverDetails = useAppSelector(selectServerDetails);
    const [showServerAlert, setShowServerAlert] = useState<boolean>(false);
    const [serverFormOpen, setServerFormOpen] = useState(false);
    const serverSaveStatus = useAppSelector(selectServerSaveStatus);
    const [pinging, setPinging] = useState<boolean>(false);
    const navigate = useNavigate();
    const { pathname } = useLocation();

    const isLoading = (serverSelectStatus === FetchStatus.Pending || serverDetailsSelectStatus === FetchStatus.Pending) && !pinging;

    useEffect(() => {
        if (authToken && serverSelectStatus !== FetchStatus.Pending) {
            dispatch(fetchServer({ guid: guid as string, token: authToken }));
        }
    }, [authToken, guid]);

    useEffect(() => {
        setShowServerAlert(serverDetailsSelectStatus === FetchStatus.Failed);
    }, [serverDetailsSelectStatus]);

    useEffect(() => {
        if (authToken && serverSelectStatus === FetchStatus.Succeeded && serverDetailsSelectStatus !== FetchStatus.Pending) {
            dispatch(fetchServerDetails({ guid: guid as string, token: authToken }));
        }
    }, [serverSelectStatus, authToken, guid, server]);

    const onTokenCopy = () => toast.success('Copied', {
        autoClose: 1000,
    });

    const handleServerFormClose = (): void => {
        setServerFormOpen(false);
    };

    const handleServerFormOpen = (): void => {
        setServerFormOpen(true);
    };

    const handleSubmit = async (server: ServerSavePayload) => {
        dispatch(editServer(server));
    };

    const handlePing = async (guid: string, authToken: string) => {
        setPinging(true);
        await dispatch(fetchServerDetails({ guid: guid as string, token: authToken }));
        setPinging(false);
    };

    const handleDomainSelect = (domain: Domain): void => {
        dispatch(domainFetched(domain));
        navigate(`${pathname}/domain/${encode(domain.servername)}`);
    }

    const uptime = serverDetails?.uptime ? moment.duration(serverDetails.uptime, 'seconds').humanize() : emptyPlaceholder;

    if (!isLoading && !server) {
        return <Error404 />
    }

    return (
        !isLoading ?
            <>
                <Breadcrumb pageName='Server' />
                <div className='flex flex-col gap-9'>
                    <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                        {showServerAlert && (
                            <Alert color="failure" onDismiss={() => setShowServerAlert(false)} className='mb-4 rounded-sm'>
                                Failed to connect to the server agent. Please check server IP address, port and token provided.
                            </Alert>
                        )}
                        <div className="border-b border-stroke px-4 py-5 dark:border-strokedark md:px-6 xl:px-9">
                            <div className="items-center sm:flex">
                                <div className="w-full flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
                                    <Avatar img={getOsIcon(server?.os_code || '')} size='lg' rounded className='justify-start'>
                                        <div className="space-y-2 font-medium dark:text-white">
                                            <h3 className="inline-block text-2xl font-medium text-black hover:text-primary dark:text-white">{server?.name || emptyPlaceholder}</h3>
                                            <div>
                                                {
                                                    serverDetails?.is_active
                                                        ? <Badge color='success' size='sm' className='inline'>Active</Badge>
                                                        : <Badge color='failure' size='sm' className='inline'>Inactive</Badge>
                                                }
                                            </div>
                                        </div>
                                    </Avatar>
                                    <div className='flex gap-3'>
                                        <Button color='blue' onClick={() => handleServerFormOpen()}>
                                            <HiMiniPencil className="mr-2 h-5 w-4" />
                                            Edit
                                        </Button>
                                        {authToken && (
                                            <Button color='success' size='sm' onClick={() => handlePing(guid as string, authToken)} isProcessing={pinging}>
                                                {!pinging && <HiMiniSignal className="mr-2 h-5 w-5" />}
                                                Ping
                                            </Button>
                                        )}

                                    </div>
                                </div>
                            </div>
                        </div>
                        <div className="p-4 md:p-6 xl:p-9">
                            <div className='flex flex-col gap-3'>
                                <div className='flex gap-6 flex-col sm:flex-row sm:gap-3'>
                                    <div className='sm:basis-1/2'>
                                        <h4 className="mb-3 font-bold text-black dark:text-white uppercase">Server Information</h4>
                                        <dl className='flex flex-col gap-3'>
                                            <div>
                                                <dt>Operating System</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    {getOsName(serverDetails?.os_code || server?.os_code || '')}{serverDetails ? ` ${serverDetails.os_version}` : server ? ` ${server.os_version}` : ''}
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Kernal</dt>
                                                <dd className="font-bold text-black dark:text-white">
                                                    {serverDetails ? `${serverDetails.kernal_version} ${serverDetails.kernal_arch}` : emptyPlaceholder}
                                                </dd>
                                            </div>
                                            <div>
                                                <dt>Hostname</dt>
                                                <dd className="font-bold text-black dark:text-white">{serverDetails?.hostname || emptyPlaceholder}</dd>
                                            </div>
                                            <div>
                                                <dt>IPv4 Address</dt>
                                                <dd className="font-bold text-black dark:text-white">{server?.ipv4_address || emptyPlaceholder}</dd>
                                            </div>
                                            <div>
                                                <dt>IPv6 Address</dt>
                                                <dd className="font-bold text-black dark:text-white">{server?.ipv6_address || emptyPlaceholder}</dd>
                                            </div>
                                            <div>
                                                <dt>Uptime</dt>
                                                <dd className="font-bold text-black dark:text-white">{uptime}</dd>
                                            </div>
                                        </dl>
                                    </div>
                                    <div className='sm:basis-1/2'>
                                        <h4 className="mb-3 font-bold text-black dark:text-white uppercase">Agent Information</h4>
                                        <dl className='flex flex-col gap-3'>
                                            <div>
                                                <dt>Agent Version</dt>
                                                <dd className="font-bold text-black dark:text-white">{serverDetails?.agent_version || server?.agent_version || emptyPlaceholder}</dd>
                                            </div>
                                            <div>
                                                <dt>Agent Port</dt>
                                                <dd className="font-bold text-black dark:text-white">{server?.agent_port || emptyPlaceholder}</dd>
                                            </div>
                                            <div>
                                                <dt>Token</dt>
                                                <dd className="font-bold text-black dark:text-white flex gap-3 items-baseline">
                                                    <span>********************************</span>
                                                    <Popover placement='left' trigger='click' content={
                                                        <div className="text-sm text-gray-500 dark:text-gray-400">
                                                            <div className="px-3 py-2">
                                                                <p>{server?.token || emptyPlaceholder}</p>
                                                            </div>
                                                        </div>
                                                    }>
                                                        <div className='cursor-pointer'>
                                                            <Tooltip content='show'>
                                                                <HiMiniEye />
                                                            </Tooltip>
                                                        </div>
                                                    </Popover>
                                                    <Tooltip content='copy'>
                                                        <CopyToClipboard text={server?.token || ''} onCopy={onTokenCopy}>
                                                            <HiMiniClipboard className='cursor-pointer' />
                                                        </CopyToClipboard>
                                                    </Tooltip>
                                                </dd>
                                            </div>
                                        </dl>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <ServerDomainList domains={serverDetails?.domains || []} onDomainClick={handleDomainSelect} />
                    <ServerEditDrawer
                        open={serverFormOpen}
                        authToken={authToken || ''}
                        loading={serverSaveStatus === FetchStatus.Pending}
                        onSubmit={handleSubmit}
                        onClose={handleServerFormClose}
                        server={server}
                    />
                </div >
            </> : <Loader />
    );
}

export default Server;
