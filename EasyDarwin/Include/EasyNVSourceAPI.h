/*
	Copyright (c) 2013-2015 EasyDarwin.ORG.  All rights reserved.
	Github: https://github.com/EasyDarwin
	WEChat: EasyDarwin
	Website: http://www.easydarwin.org
*/
#ifndef _Easy_NVS_API_H
#define _Easy_NVS_API_H

#define WIN32_LEAN_AND_MEAN
//#include <windows.h>
#include <winsock2.h>

#ifdef _WIN32
#define EasyNVS_API  __declspec(dllexport)
#define Easy_APICALL  __stdcall
#else
#define EasyNVS_API
#define Easy_APICALL 
#endif

#define Easy_NVS_Handle void*

//ý������
#ifndef MEDIA_TYPE_VIDEO
#define MEDIA_TYPE_VIDEO		0x00000001
#endif
#ifndef MEDIA_TYPE_AUDIO
#define MEDIA_TYPE_AUDIO		0x00000002
#endif
#ifndef MEDIA_TYPE_EVENT
#define MEDIA_TYPE_EVENT		0x00000004
#endif
#ifndef MEDIA_TYPE_RTP
#define MEDIA_TYPE_RTP			0x00000008
#endif
#ifndef MEDIA_TYPE_SDP
#define MEDIA_TYPE_SDP			0x00000010
#endif

//video codec
#define	VIDEO_CODEC_H264	0x1C
#define	VIDEO_CODEC_MJPEG	0x08
#define	VIDEO_CODEC_MPEG4	0x0D
//audio codec
#define AUDIO_CODEC_MP4A	0x15002		//86018
#define AUDIO_CODEC_PCMU	0x10006		//65542


//��������
typedef enum __RTP_CONNECT_TYPE
{
	RTP_OVER_TCP	=	0x01,
	RTP_OVER_UDP
}RTP_CONNECT_TYPE;

//֡����
#ifndef FRAMETYPE_I
#define FRAMETYPE_I		0x01
#endif
#ifndef FRAMETYPE_P
#define FRAMETYPE_P		0x02
#endif
#ifndef FRAMETYPE_B
#define FRAMETYPE_B		0x03
#endif

//֡��Ϣ
typedef struct 
{
	unsigned int	codec;			//�����ʽ
	unsigned char	type;			//֡����
	unsigned char	fps;			//֡��
	unsigned char	reserved1;
	unsigned char	reserved2;

	unsigned short	width;			//��
	unsigned short  height;			//��
	unsigned int	sample_rate;	//������
	unsigned int	channels;		//����
	unsigned int	length;			//֡��С
	unsigned int    rtptimestamp;	//rtp timestamp
	unsigned int	timestamp_sec;	//��
	
	float			bitrate;
	float			losspacket;
}NVS_FRAME_INFO;

/*
//�ص�:
_mediatype:		MEDIA_TYPE_VIDEO	MEDIA_TYPE_AUDIO	MEDIA_TYPE_EVENT	
�����EasyNVS_OpenStream�еĲ���outRtpPacket��Ϊ1, ��ص��е�_mediatypeΪMEDIA_TYPE_RTP, pbufΪ���յ���RTP��(����rtpͷ��Ϣ), frameinfo->lengthΪ����
*/
typedef int (CALLBACK *NVSourceCallBack)( int _chid, int *_chPtr, int _mediatype, char *pbuf, NVS_FRAME_INFO *frameinfo);

extern "C"
{
	//��ȡ�������
	EasyNVS_API int Easy_APICALL EasyNVS_GetErrCode();

	EasyNVS_API int Easy_APICALL EasyNVS_Init(Easy_NVS_Handle *handle);
	EasyNVS_API int Easy_APICALL EasyNVS_Deinit(Easy_NVS_Handle *handle);

	EasyNVS_API int Easy_APICALL EasyNVS_SetCallback(Easy_NVS_Handle handle, NVSourceCallBack _callback);

	EasyNVS_API int Easy_APICALL EasyNVS_OpenStream(Easy_NVS_Handle handle, int _channelid, char *_url, RTP_CONNECT_TYPE _connType, unsigned int _mediaType, char *_username, char *_password, void *userPtr, int _reconn/*1000��ʾ������,���������Ͽ��Զ�����, ����ֵΪ���Ӵ���*/, int outRtpPacket/*Ĭ��Ϊ0,���ص����������֡, ���Ϊ1,�����RTP��*/);
	EasyNVS_API int Easy_APICALL EasyNVS_CloseStream(Easy_NVS_Handle handle);
};





#endif