import { Badge, Button, Popover, Spinner } from 'flowbite-react';
import { FC, useState } from 'react';
import { HiArrowPath, HiChevronDoubleDown, HiMiniEye } from 'react-icons/hi2';
import { RenewalLog } from '../types';
import moment from 'moment';
import empty from '../../../images/empty.png';

type RenewalLogsProps = {
    logs: RenewalLog[]
    loading: boolean
    onRefresh: () => void;
};

const expandLimit = 5;

const RenewalLogs: FC<RenewalLogsProps> = ({ logs, onRefresh, loading }) => {
    const [expandedCount, setExpandedCount] = useState<number>(expandLimit);
    let slicedLogs = logs.slice(0, expandedCount);

    const expandLogs = () => {
        const limit = expandedCount + expandLimit

        setExpandedCount(limit);
        slicedLogs = logs.slice(0, limit);
    };

    const getLogContent = (message: string) => {
        return (
            <div className="text-sm text-gray-500 dark:text-gray-400 max-w-[300px]">
                <div className="border-b border-gray-200 bg-gray-100 px-3 py-2 dark:border-gray-600 dark:bg-gray-700">
                    <h3 id="default-popover" className="font-semibold text-gray-900 dark:text-white">
                        Message
                    </h3>
                </div>
                <div className="px-3 py-2">
                    <span dangerouslySetInnerHTML={{ __html: message }} className='break-words' />
                </div>
            </div>
        )
    }

    return (
        <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
            <div className="p-4 md:p-6 xl:p-7.5">
                <div className="flex items-start justify-between">
                    <div className='flex justify-between items-center w-full'>
                        <h2 className="text-xl font-bold text-black dark:text-white">Latest Auto Renewal Logs</h2>
                        <Button color='blue' onClick={onRefresh}>
                            {loading ? <Spinner /> : (
                                <>
                                    <HiArrowPath className="mr-2 h-5 w-4" />
                                    Refresh
                                </>
                            )}
                        </Button>
                    </div>
                </div>
            </div>
            <div className='border-b border-stroke px-4 dark:border-strokedark md:px-6 xl:px-7.5' />
            <div className='p-4 md:p-6 xl:p-7.5'>
                <div className="max-w-full overflow-x-auto">
                    <div className="mx-auto max-w-[80px]">
                        {!logs.length && <img src={empty} />}
                    </div>
                    {logs.length !== 0 && (
                        <>
                            <table className="w-full table-auto">
                                <thead>
                                    <tr className="bg-gray-2 text-left dark:bg-meta-4">
                                        <th className="min-w-[150px] py-4 px-4 font-medium text-black dark:text-white xl:pl-11">
                                            Domain
                                        </th>
                                        <th className="min-w-[120px] py-4 px-4 font-medium text-black dark:text-white">
                                            Date
                                        </th>
                                        <th className="min-w-[100px] py-4 px-4 font-medium text-black dark:text-white">
                                            Status
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {slicedLogs.map((log, key) => (
                                        <tr key={key}>
                                            <td className="border-b border-[#eee] py-5 px-4 pl-9 dark:border-strokedark xl:pl-11">
                                                <h5 className="font-medium text-black dark:text-white">
                                                    {log.domainName}
                                                </h5>
                                                <p className="text-sm">{log.serverName}</p>
                                            </td>
                                            <td className="border-b border-[#eee] py-5 px-4 dark:border-strokedark">
                                                {moment(log.createdAt).format('LL, h:mm')}
                                            </td>
                                            <td className='border-b border-[#eee] py-5 px-4 dark:border-strokedark'>
                                                {
                                                    log.message === ''
                                                        ? <Badge color='success' size='sm' className='inline'>Success</Badge>
                                                        : (
                                                            <div className='flex gap-2 items-center'>
                                                                <Badge color='failure' size='sm' className='inline'>Error</Badge>
                                                                <Popover content={getLogContent(log.message)} placement='left'>
                                                                    <span className='cursor-pointer'><HiMiniEye size={20} className='hover:text-blue-500 mr-2' /></span>
                                                                </Popover>
                                                            </div>
                                                        )
                                                }
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                            {
                                slicedLogs.length >= logs.length ? null : (
                                    <div className='py-4'>
                                        <Button size='sm' color='dark' className='mx-auto' onClick={expandLogs}>
                                            <HiChevronDoubleDown className='h-6 w-6' />
                                        </Button>
                                    </div>
                                )
                            }
                        </>
                    )}

                </div>
            </div>
        </div >
    );
};

export default RenewalLogs;
