import Breadcrumb from '../../components/Breadcrumb';
import { Button } from 'flowbite-react';
import CertificateItem from '../../features/certificate/components/CertificateItem';
import useAuthToken from '../../features/auth/hooks';
import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { fetchCertificates, selectCertificates, selectCertificatesFetchStatus } from '../../features/certificate/certificatesSlice';
import { useEffect } from 'react';
import { FetchStatus } from '../../app/types';
import { useParams } from 'react-router-dom';
import Loader from '../../components/Loader/Loader';
import empty from '../../images/empty.png';
import { downloadCertificateApi } from '../../features/certificate/certificateApi';
import { toast } from 'react-toastify';

const CertificateList = () => {
    const { guid } = useParams();
    const [authToken] = useAuthToken();
    const dispatch = useAppDispatch();
    const certificatesSelectStatus = useAppSelector(selectCertificatesFetchStatus);
    const certificates = useAppSelector(selectCertificates);

    const isLoading = certificatesSelectStatus === FetchStatus.Pending;
    const certNames = Object.keys(certificates);

    useEffect(() => {
        if (authToken && certificatesSelectStatus !== FetchStatus.Pending) {
            dispatch(fetchCertificates({
                guid: guid as string,
                token: authToken,
            }));
        }
    }, [authToken]);

    const handleCertificateDownload = async (name: string): Promise<any> => {
        if (!authToken || !guid) {
            return;
        }

        try {
            const response = await downloadCertificateApi({
                guid,
                name,
                token: authToken,
            });

            const url = window.URL.createObjectURL(new Blob([response.content]));
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', response.name);
            document.body.appendChild(link);
            link.click();
            link.remove();
        } catch (err) {
            const error = err as Error;
            toast.error(error.message);
        }
    };

    return (
        isLoading
            ? <Loader />
            : <>
                <Breadcrumb pageName='Certificates'>
                    <Button color='blue'>Add</Button>
                </Breadcrumb>
                <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                    <div className="border-b border-stroke px-4 py-5 dark:border-strokedark md:px-6 xl:px-7.5">
                        <div className="flex items-center gap-2">
                            <div className="w-3/12 md:w-4/12">
                                <span className="font-medium">Vendor</span>
                            </div>
                            <div className="w-4/12">
                                <span className="font-medium">Names</span>
                            </div>
                            <div className="w-3/12">
                                <span className="font-medium">Expires</span>
                            </div>
                            <div className="w-2/12 md:w-1/12 text-center">
                                <span className="font-medium">Actions</span>
                            </div>
                        </div>
                    </div>
                    <div className="p-4 md:p-6 xl:p-7.5">
                        <div className="mx-auto max-w-[80px]">
                            {!certNames.length && <img src={empty} />}
                        </div>
                        {
                            Object.entries(certificates).map(([name, certificate]) => (
                                <CertificateItem
                                    name={name}
                                    certificate={certificate}
                                    key={certificate.cn}
                                    onCertificateDownload={handleCertificateDownload}
                                />
                            ))
                        }
                    </div>
                </div>
            </>
    );
}

export default CertificateList;
