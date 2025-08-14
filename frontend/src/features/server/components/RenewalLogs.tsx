import { Badge, Button } from 'flowbite-react';
import { FC } from 'react';
import { HiArrowPath } from 'react-icons/hi2';
import { RenewalLog } from '../types';

const emptyPlaceholder = '----------';

type RenewalLogsProps = {
};

const items: RenewalLog[] = [
    {
        domainName: 'example.com',
        serverName: 'Test Server',
        message: 'Some Error',
        date: '2025-03-34 12:00',
    },
    {
        domainName: 'example2.com',
        serverName: 'Test Server2',
        message: '',
        date: '2025-03-34 13:00',
    },
];

const RenewalLogs: FC<RenewalLogsProps> = () => {
    return (
        <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
            <div className="p-4 md:p-6 xl:p-7.5">
                <div className="flex items-start justify-between">
                    <div className='flex justify-between items-center w-full'>
                        <h2 className="text-xl font-bold text-black dark:text-white">Renewal Logs</h2>
                        <Button color='blue'>
                            <HiArrowPath className="mr-2 h-5 w-4" />
                            Refresh
                        </Button>
                    </div>
                </div>
            </div>
            <div className='border-b border-stroke px-4 dark:border-strokedark md:px-6 xl:px-7.5' />
            <div className='p-4 md:p-6 xl:p-7.5'>
                <div className="max-w-full overflow-x-auto">
                    <table className="w-full table-auto">
                        <thead>
                            <tr className="bg-gray-2 text-left dark:bg-meta-4">
                                <th className="min-w-[150px] py-4 px-4 font-medium text-black dark:text-white xl:pl-11">
                                    Domain
                                </th>
                                <th className="min-w-[120px] py-4 px-4 font-medium text-black dark:text-white">
                                    Status
                                </th>
                                <th className="py-4 px-4 font-medium text-black dark:text-white">
                                    Date
                                </th>
                                <th className="min-w-[220px] py-4 px-4 font-medium text-black dark:text-white">
                                    Message
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {items.map((item, key) => (
                                <tr key={key}>
                                    <td className="border-b border-[#eee] py-5 px-4 pl-9 dark:border-strokedark xl:pl-11">
                                        <h5 className="font-medium text-black dark:text-white">
                                            {item.domainName}
                                        </h5>
                                        <p className="text-sm">{item.serverName}</p>
                                    </td>
                                    <td className='border-b border-[#eee] py-5 px-4 dark:border-strokedark'>
                                        {
                                            item.message === ''
                                                ? <Badge color='success' size='sm' className='inline'>Success</Badge>
                                                : <Badge color='failure' size='sm' className='inline'>Error</Badge>
                                        }
                                    </td>
                                    <td className="border-b border-[#eee] py-5 px-4 dark:border-strokedark">
                                        {item.date}
                                    </td>
                                    <td className="border-b border-[#eee] py-5 px-4 dark:border-strokedark">
                                        <p className="text-black dark:text-white">
                                            {item.message || emptyPlaceholder}
                                        </p>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div >
    );
};

export default RenewalLogs;
