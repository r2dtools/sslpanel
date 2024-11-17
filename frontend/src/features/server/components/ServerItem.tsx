import { Link } from 'react-router-dom';
import { Server } from '../types';
import { getOsIcon, getOsName } from '../utils';
import { Tooltip } from 'flowbite-react';
import { HiMiniTrash } from 'react-icons/hi2';

interface ServerItemProps {
    server: Server;
};

const ServerItem: React.FC<ServerItemProps> = ({ server }: ServerItemProps) => {
    console.log(server);

    return (
        <Link to={`/servers/${server.guid}`}>
            <div className="p-3 flex items-center gap-3 hover:bg-[#F8FAFD] dark:hover:bg-meta-4 hover:rounded">
                <div className="md:w-3/12 w-6/12 lg:w-4/12">
                    <div className="flex items-center gap-4">
                        <div className="2xsm:h-11 2xsm:w-full 2xsm:max-w-11 2xsm:rounded-full">
                            <img src={getOsIcon(server.os_code)} />
                        </div>
                        <div>
                            <span className="font-bold">{server.name}</span>
                            <span className="mt-1 block text-sm">{getOsName(server.os_code)} {server.os_version}</span>
                        </div>
                    </div>
                </div>
                <div className="hidden md:block md:w-7/12 lg:w-6/12 xl:w-5/12 text-sm">
                    {server.ipv4_address && <div className="font-medium">{server.ipv4_address}</div>}
                    {server.ipv6_address && <div className="font-medium">{server.ipv6_address}</div>}
                </div>
                <div className="hidden xl:block lg:w-1/12">
                    <span className="font-medium">{`v${server.agent_version}`}</span>
                </div>
                <div className="w-3/12 md:w-1/12">
                    {server.is_active ? (
                        <span className="inline-block rounded bg-green-500 dark:bg-green-500/[0.5] px-2.5 py-0.5 text-sm font-medium text-white">Active</span>
                    ) : (
                        <span className="inline-block rounded bg-red-500 dark:bg-red-500/[0.5] px-2.5 py-0.5 text-sm font-medium text-white">Inactive</span>
                    )}
                </div>
                <div className="w-3/12 md:w-1/12">
                    <button className="mx-auto block hover:text-meta-1">
                        <Tooltip content="Delete">
                            <HiMiniTrash size={20} />
                        </Tooltip>
                    </button>

                </div>
            </div>
        </Link>
    );
};

export default ServerItem;
