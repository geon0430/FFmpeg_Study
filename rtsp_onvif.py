from onvif2 import ONVIFCamera
import sys
sys.path.append("../")
from utils import custom_logger

def get_onvif_rtsp_address_list(ip, port=80, id=None, passwd=None):
    mycam = ONVIFCamera(ip, port, id, passwd, wsdl_dir='./wsdl/', adjust_time=True)
    media_service2 = mycam.create_media2_service()
    profiles = media_service2.GetProfiles()
    configurations = media_service2.GetVideoEncoderConfigurations()

    configuration_list = []
    for configuration in configurations:
        if configuration.UseCount != 0:
            configuration_list.append(configuration)

    results = []
    for p, c in zip(profiles, configuration_list):
        o = media_service2.create_type('GetStreamUri')
        o.ProfileToken = p.token
        o.Protocol = 'RTSP'
        uri_response = media_service2.GetStreamUri(o)
        
        custom_logger.debug(f"URI response: {uri_response}")

        uri = uri_response if isinstance(uri_response, str) else uri_response.Uri
        
        if passwd is None:
            rtsp_uri = uri
        else:
            rtsp_uri = f'rtsp://{id}:{passwd}@{uri[7:]}'

        result = {
            'width': c.Resolution.Width,
            'height': c.Resolution.Height,
            'codec': c.Encoding.lower(),
            'fps': int(c.RateControl.FrameRateLimit),
            'rtsp': rtsp_uri
        }

        results.append(result)
        custom_logger.debug(f"onvif | result : {result}")

    return results


def get_onvif_rtsp_address(ip, port=80, id=None, passwd=None):
    mycam = ONVIFCamera(ip, port, id, passwd, wsdl_dir='./wsdl/', adjust_time=True)
    media_service = mycam.create_media_service()
    media_service2 = mycam.create_media2_service()
    profile = media_service.GetProfiles()[1]

    o = media_service2.create_type('GetStreamUri')
    o.ProfileToken = profile.token
    o.Protocol = 'RTSP'
    uri_response = media_service2.GetStreamUri(o)

    uri = uri_response  
    rtsp_uri = uri if passwd is None else f'rtsp://{id}:{passwd}@{uri[7:]}'

    result = {
        'width': profile.VideoEncoderConfiguration.Resolution.Width,
        'height': profile.VideoEncoderConfiguration.Resolution.Height,
        'codec': profile.VideoEncoderConfiguration.Encoding.lower(),
        'fps': profile.VideoEncoderConfiguration.RateControl.FrameRateLimit,
        'rtsp': rtsp_uri
    }

    return result


if __name__ == "__main__":
    ip = ""
    id = ""
    password = ""
    profile_info = get_onvif_rtsp_address(ip, 80, id, password)
    print(profile_info)
