import { Link } from 'react-router-dom';
import { Server } from '../types';
import { getOsIcon, getOsName } from '../utils';
import { Spinner, Tooltip } from 'flowbite-react';
import { HiMiniTrash, HiMiniPencil } from 'react-icons/hi2';
import { useState } from 'react';

interface ServerItemProps {
    server: Server;
    onDelete: (id: number) => Promise<any>;
    onEdit: (server: Server) => void;
};

const ServerItem: React.FC<ServerItemProps> = ({ server, onDelete, onEdit }: ServerItemProps) => {
    const [actionLoading, setActionLoading] = useState<boolean>(false);

    const handleDelete = async (event: React.MouseEvent<SVGElement, MouseEvent>) => {
        event.preventDefault();
        event.stopPropagation();

        setActionLoading(true);
        await onDelete(server.id);
        setActionLoading(false);
    };

    const handleEdit = (event: React.MouseEvent<SVGElement, MouseEvent>) => {
        event.preventDefault();
        event.stopPropagation();

        onEdit(server);
    };

    const preventClick = (event: React.MouseEvent) => {
        event.preventDefault();
        event.stopPropagation();
    };

    return (
        <Link to={`/servers/${server.guid}`}>
            <div className="p-3 flex items-center gap-3 hover:bg-[#F8FAFD] dark:hover:bg-meta-4 hover:rounded">
                <div className="md:w-3/12 w-6/12 lg:w-4/12">
                    <div className="flex items-center gap-4">
                        <div className="2xsm:h-11 2xsm:w-full 2xsm:max-w-11 2xsm:rounded-full">
                            <img src={getOsIcon(server.os_code)} />
                        </div>
                        <div>
                            <span className="font-bold text-black dark:text-white">{server.name}</span>
                            <span className="mt-1 block text-sm">{getOsName(server.os_code)} {server.os_version}</span>
                        </div>
                    </div>
                </div>
                <div className="hidden md:block md:w-7/12 lg:w-6/12 xl:w-5/12 text-sm text-black dark:text-white">
                    {server.ipv4_address && <div className="font-medium">{server.ipv4_address}</div>}
                    {server.ipv6_address && <div className="font-medium">{server.ipv6_address}</div>}
                </div>
                <div className="hidden xl:block lg:w-1/12">
                    <span className="font-medium">{server.agent_version ? server.agent_version : ''}</span>
                </div>
                <div className="w-3/12 md:w-1/12">
                    {server.is_active ? (
                        <span className="inline-block rounded bg-green-500 dark:bg-green-500 px-2.5 py-0.5 text-sm font-medium text-white">Active</span>
                    ) : (
                        <span className="inline-block rounded bg-red-500 dark:bg-red-500 px-2.5 py-0.5 text-sm font-medium text-white">Inactive</span>
                    )}
                </div>
                <div className="w-3/12 md:w-1/12 text-center" onClick={preventClick}>
                    {
                        actionLoading
                            ? <Spinner />
                            : (
                                <button className="flex justify-between mx-auto block ">
                                    <Tooltip content="Edit">
                                        <HiMiniPencil size={20} className='hover:text-blue-500 mr-2' onClick={handleEdit} />
                                    </Tooltip>
                                    <Tooltip content="Delete" >
                                        <HiMiniTrash size={20} className='hover:text-red-500' onClick={handleDelete} />
                                    </Tooltip>
                                </button>
                            )
                    }
                </div>
            </div>
        </Link>
    );
};

export default ServerItem;
