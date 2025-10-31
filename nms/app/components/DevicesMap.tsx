import pulsarCaseImg from '../assets/pulsar_case.png';
import pulsarCaseRed from '../assets/pulsar_case_red.png';
import pulsarCaseBlue from '../assets/pulsar_case_blue.png';
import flyerCartRed from '../assets/flyer_cart_red.png';
import flyerCartBlue from '../assets/flyer_cart_blue.png';
import flyerCartImg from '../assets/flyer_cart.png';

import React, { useState, useEffect, useMemo } from "react";
import CardTitleRow from "./layout/CardTitleRow";
import { Map as PigeonMap, Overlay } from 'pigeon-maps';
import Typography from '@mui/material/Typography';
import MapIcon from '@mui/icons-material/Map';
import CancelIcon from '@mui/icons-material/Cancel';
import Fade from '@mui/material/Fade';
import CloseIcon from '@mui/icons-material/Close';
import Tooltip from '@mui/material/Tooltip';
import { withStyles } from '@mui/styles';

import PublicIcon from '@mui/icons-material/Public'; 
import WifiIcon from '@mui/icons-material/Wifi'; 
import NetworkCellIcon from '@mui/icons-material/NetworkCell'; 
import CheckCircleIcon from '@mui/icons-material/CheckCircle'; 
import { Box, Card, CardContent, CardHeader, Grid, IconButton } from "@mui/material";
import { useNavigate, useParams } from 'react-router-dom';
import MagmaAPI from '../api/MagmaAPI';

export const DevicesMap = () => {
  const HtmlTooltip = withStyles((theme: any) => ({
    tooltip: {
      backgroundColor: '#f5f5f9',
      color: 'rgba(0, 0, 0, 0.87)',
      maxWidth: 220,
      fontSize: theme.typography.pxToRem(12),
      border: '1px solid #dadde9',
    },
  }))(Tooltip);

  interface DeviceT{
    imei: string;
    lat: number;
    lng: number;
    backhaul: boolean;
    private_5g: boolean;
    data_on_5g: boolean;
    type: string;
    name:string;
    bytes_received?: number;
    bytes_sent?: number;
    subscribers?: number;
    device_type?: string;
    plotLat?: number;
    plotLng?: number;
  }

  const imagesSwitch: Record<string, Record<'red'|'blue', string>> = {
    "Pulsar Case":{
      "red": pulsarCaseRed,
      "blue": pulsarCaseBlue,
    },
    "Flyer Cart":{
      "red": flyerCartRed,
      "blue": flyerCartBlue,
    }
  };

  const {networkId} = useParams();
  const navigate = useNavigate();

  const [actualMarker, setActualMarker] = useState<DeviceT | undefined>(undefined);
  const [devices, setDevices] = useState<DeviceT[]>([]);
  const [mapZoom, setMapZoom] = useState(4);

  const verifyAllActive = (marker: DeviceT) => {
    return marker.backhaul && marker.private_5g && marker.data_on_5g;
  };

  const getPrometheusData = async () => {
    const res = await MagmaAPI.metrics.networksNetworkIdPrometheusQueryGet({
      networkId: networkId!,
      query: `device_network_service`,
    });
    return res;
  };

  const getData5gValue = async (imei: string) =>{
    const data5gValue = await MagmaAPI.metrics.networksNetworkIdPrometheusQueryGet({
      networkId: networkId!,
      query: `device_private5g_active{imei="${imei}"}`,
    });
    const dataOn5G =
      data5gValue.data.data.result.length > 0
        ? parseFloat(data5gValue.data.data.result[0].value[1])
        : 0;
    return dataOn5G;
  };

  const getBackhaulStatus = async (imei: string) =>{
    const backhaulRes = await MagmaAPI.metrics.networksNetworkIdPrometheusQueryGet({
      networkId: networkId!,
      query: `device_comm_status{imei="${imei}"}`,
    });
    const backhaulValue =
      backhaulRes.data.data.result.length > 0
        ? parseFloat(backhaulRes.data.data.result[0].value[1])
        : 0;
    return backhaulValue;
  };

  const getBytesRecieved = async (imei: string) =>{
    const brecievedRes = await MagmaAPI.metrics.networksNetworkIdPrometheusQueryGet({
      networkId: networkId!,
      query: `device_bytes_received{imei="${imei}"}`,
    });
    const brecievedValue =
      brecievedRes.data.data.result.length > 0
        ? parseFloat(brecievedRes.data.data.result[0].value[1])
        : 0;
    return brecievedValue;
  };

  const getBytesSent = async (imei: string) =>{
    const bsentRes = await MagmaAPI.metrics.networksNetworkIdPrometheusQueryGet({
      networkId: networkId!,
      query: `device_bytes_sent{imei="${imei}"}`,
    });
    const bsentValue =
      bsentRes.data.data.result.length > 0
        ? parseFloat(bsentRes.data.data.result[0].value[1])
        : 0;
    return bsentValue;
  };

  const getLatitude = async (imei: string) =>{
    const latitudeRes = await MagmaAPI.metrics.networksNetworkIdPrometheusQueryGet({
      networkId: networkId!,
      query: `device_latitude{imei="${imei}"}`,
    });
    const latitudeValue =
      latitudeRes.data.data.result.length > 0
        ? parseFloat(latitudeRes.data.data.result[0].value[1])
        : 0;
    return latitudeValue;
  };

  const getLongitude = async (imei: string) =>{
    const longitudeRes = await MagmaAPI.metrics.networksNetworkIdPrometheusQueryGet({
      networkId: networkId!,
      query: `device_longitude{imei="${imei}"}`,
    });
    const longitudeValue =
      longitudeRes.data.data.result.length > 0
        ? parseFloat(longitudeRes.data.data.result[0].value[1])
        : 0;
    return longitudeValue;
  };

  // Função de conversão de metros -> graus (offset fixo, independente do zoom)
  function metersToDegreeOffsets(meters: number, latDeg: number) {
    const latRad = (latDeg * Math.PI) / 180;
    const degPerMeterLat = 1 / 111320;
    const degPerMeterLng = 1 / (111320 * Math.cos(latRad));
    return {
      dLatDeg: meters * degPerMeterLat,
      dLngDeg: meters * degPerMeterLng,
    };
  }

  // Aplica offset SOMENTE quando o zoom está alto; offset é fixo em metros
  function applyOffsetForCoincident(devs: DeviceT[], zoom: number): DeviceT[] {
    const OFFSET_ZOOM_THRESHOLD = 12; // começa a separar a partir deste zoom
    const METERS_RADIUS = 20;         // raio de separação fixo (~20 m)

    // Longe: sem offset
    if (zoom < OFFSET_ZOOM_THRESHOLD) {
      return devs.map(d => ({ ...d, plotLat: d.lat, plotLng: d.lng }));
    }

    const keyOf = (d: DeviceT) => `${d.lat.toFixed(6)}|${d.lng.toFixed(6)}`;
    const groups = new Map<string, number[]>();

    devs.forEach((d, idx) => {
      const k = keyOf(d);
      if (!groups.has(k)) groups.set(k, []);
      groups.get(k)!.push(idx);
    });

    const out = devs.map(d => ({ ...d, plotLat: d.lat, plotLng: d.lng }));

    groups.forEach(indices => {
      if (indices.length <= 1) return;

      const n = indices.length;
      const angleStep = (2 * Math.PI) / n;

      indices.forEach((ix, i) => {
        const d = devs[ix];

        // deslocamento fixo convertido para graus na latitude local
        const { dLatDeg, dLngDeg } = metersToDegreeOffsets(METERS_RADIUS, d.lat);

        const angle = i * angleStep;
        const oLat = dLatDeg * Math.cos(angle);
        const oLng = dLngDeg * Math.sin(angle);

        out[ix].plotLat = d.lat + oLat;
        out[ix].plotLng = d.lng + oLng;
      });
    });

    return out;
  }

  const plottedDevices = useMemo(
    () => applyOffsetForCoincident(devices, mapZoom),
    [devices, mapZoom]
  );

  useEffect(() => {
    let cancelled = false;
    let timer: number | undefined;

    const fetchData = async () => {
      try {
        const response = await getPrometheusData();
        let device: DeviceT[] = await Promise.all(
          response.data.data.result.map(async (item: any) => {
            const private5gValue = item.value ? parseFloat(item.value[1]) : 0;
            const data5gValue = await getData5gValue(item.metric.imei);
            const backhaulValue = await getBackhaulStatus(item.metric.imei);
            const bytesRecieved = await getBytesRecieved(item.metric.imei);
            const bytesSent = await getBytesSent(item.metric.imei);
            const latitude = await getLatitude(item.metric.imei);
            const longitude = await getLongitude(item.metric.imei);

            return {
              imei: item.metric.imei,
              lat: latitude ? latitude : 32.99467491596908,
              lng: longitude ? longitude : -97.05818548575185,
              backhaul: backhaulValue === 1,
              private_5g: private5gValue > 0,
              data_on_5g: private5gValue > 0 && data5gValue === 1,
              type: (item.metric.name).includes("Flyer") ? "Flyer Cart" : "Pulsar Case",
              name: item.metric.name,
              device_type: item.metric.type || undefined,
              bytes_received: bytesRecieved,
              bytes_sent: bytesSent, 
            } as DeviceT;
          })
        );
        if (!cancelled) setDevices(device);
      } catch (error) {
        if (!cancelled) console.error("Error fetching data:", error);
      } finally {
        if (!cancelled) {
          timer = window.setTimeout(fetchData, 60_000);
        }
      }
    };

    fetchData();

    return () => {
      cancelled = true;
      if (timer) window.clearTimeout(timer);
    };
  }, []);

  return(
    <>
      <CardTitleRow icon={MapIcon} label="Devices Map" />
      <PigeonMap
        defaultCenter={[36.22853, -93.03602]}
        defaultZoom={4}
        height={400}
        onBoundsChanged={({ zoom }) => setMapZoom(zoom)}
      >
        {plottedDevices.map((marker) => (
          <Overlay
            key={marker.imei}
            anchor={[marker.plotLat ?? marker.lat, marker.plotLng ?? marker.lng]}
          >
            <div
              onClick={() => { setActualMarker(marker); }}
              style={{
                cursor: 'pointer',
                transform: 'translate(-50%, -50%)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                width: 30,
                height: 30,
              }}
            >
              {verifyAllActive(marker) ? (
                <img src={imagesSwitch[marker.type]["blue"]} alt="" style={{width: 30, height: 30}}/>
              ) : (
                <img src={imagesSwitch[marker.type]["red"]} alt="" style={{width: 30, height: 30}}/>
              )}
            </div>
          </Overlay>
        ))}
      </PigeonMap>

      {actualMarker && (
        <Fade in={!!actualMarker} timeout={500}>
          <Card elevation={1} style={{marginTop: 16}}>
            <CardHeader
              avatar={
                <img
                  src={actualMarker.type === 'Pulsar Case' ? pulsarCaseImg : flyerCartImg}
                  alt="Device"
                  style={{ width: 150, height: 150 }}
                />
              }
              action={
                <IconButton onClick={() => setActualMarker(undefined)}>
                  <CloseIcon />
                </IconButton>
              }
              title="Device Details"
              subheaderTypographyProps={{ style: { fontSize:18, fontWeight: 400 } }}
              titleTypographyProps={{ style: { fontSize:18, fontWeight: 400 } }}
              subheader={
                <>
                  <div>
                    Lat: {actualMarker.lat ? (actualMarker.lat).toFixed(6) : ''} | Lng: {actualMarker.lng ? (actualMarker.lng).toFixed(6) : ''}
                  </div>
                  <div>Type: {actualMarker.type}</div>
                  <div>Type: {actualMarker.name}</div>
                  {actualMarker.device_type && <div>Device Model: {actualMarker.device_type}</div>}
                  <div>Active Subscribers: {actualMarker.subscribers ? actualMarker.subscribers : "-"}</div>
                </>
              }
            />
            
            <CardContent>
              <Grid container spacing={2}>
                {/* Backhaul */}
                <HtmlTooltip
                  placement="top-start"
                  title={
                    <React.Fragment>
                      <Typography color="inherit">Details:</Typography>
                      <p>
                        Device Comm Status: { actualMarker.backhaul ? "Connected" : "Disconnected" }
                      </p>
                    </React.Fragment>
                  }
                >
                  <Grid item xs={12} md={4}>
                    <Box display="flex" alignItems="center">
                      <PublicIcon style={{marginRight: 8, fontSize:50, color: actualMarker.backhaul ? 'green' : 'red'}}/>
                      <Typography variant="body1" style={{marginRight: 8}}>
                        Backhaul Internet-In
                      </Typography>
                      {actualMarker.backhaul ? (
                        <>
                          <CheckCircleIcon style={{marginRight: 4, color: 'green'}} />
                          <Typography variant="body2" style={{color: 'green'}}>Active</Typography>
                        </>
                      ) : (
                        <>
                          <CancelIcon style={{marginRight: 4, color: 'red'}} />
                          <Typography variant="body2" style={{color: 'red'}}>Inactive</Typography>
                        </>
                      )}
                    </Box>
                  </Grid>
                </HtmlTooltip>

                {/* Private 5G */}
                <HtmlTooltip
                  placement="top-start"
                  title={
                    <React.Fragment>
                      <Typography color="inherit">Details:</Typography>
                      <p>
                        Network Service Type: {
                          // mantém a lógica original usando == para mapear boolean para 0/1
                          actualMarker.private_5g == 0 ? "Unknown" : actualMarker.private_5g == 1 ? "4G" : actualMarker.private_5g == 2 ? "5G" : undefined
                        }
                      </p>
                    </React.Fragment>
                  }
                >
                  <Grid item xs={12} md={4}>
                    <Box display="flex" alignItems="center">
                      <WifiIcon style={{marginRight: 8, fontSize:50, color: actualMarker.private_5g == 1 ? 'orange' : actualMarker.private_5g ? 'green' : 'red'}} />
                      <Typography variant="body1" style={{marginRight: 8}}>
                        Private 5G
                      </Typography>
                      {actualMarker.private_5g ? (
                        <>
                          <CheckCircleIcon style={{marginRight: 4, color: actualMarker.private_5g == 1 ? 'orange' : 'green'}} />
                          <Typography variant="body2" style={{color: actualMarker.private_5g == 1 ? 'orange' : 'green'}}>Active</Typography>
                        </>
                      ) : (
                        <>
                          <CancelIcon style={{marginRight: 4, color: 'red'}} />
                          <Typography variant="body2" style={{color: 'red'}}>Inactive</Typography>
                        </>
                      )}
                    </Box>
                  </Grid>
                </HtmlTooltip>

                {/* 5G data */}
                <HtmlTooltip
                  placement="top-start"
                  title={
                    <React.Fragment>
                      <Typography color="inherit">Details:</Typography>
                      <p>
                        Bytes Received: {actualMarker.bytes_received ? actualMarker.bytes_received : "-"} <br/>
                        Bytes Sent: {actualMarker.bytes_sent ? actualMarker.bytes_sent : "-"}
                      </p>
                    </React.Fragment>
                  }
                >
                  <Grid item xs={12} md={4}>
                    <Box display="flex" alignItems="center">
                      <NetworkCellIcon
                        onClick={() => navigate(`/nms/${networkId}/metrics/grafana`)}
                        style={{marginRight: 8, fontSize:50, color: actualMarker.private_5g == 1 ? 'orange' : actualMarker.data_on_5g ? 'green' : 'red', cursor:'pointer'}}
                      />
                      <Typography variant="body1" style={{marginRight: 8}}>
                        Data on Private 5G
                      </Typography>
                      {actualMarker.data_on_5g ? (
                        <>
                          <CheckCircleIcon style={{marginRight: 4, color: actualMarker.private_5g == 1 ? 'orange' : 'green'}} />
                          <Typography variant="body2" style={{color: actualMarker.private_5g == 1 ? 'orange' : 'green'}}>Active</Typography>
                        </>
                      ) : (
                        <>
                          <CancelIcon style={{marginRight: 4, color: 'red'}} />
                          <Typography variant="body2" style={{color: 'red'}}>Inactive</Typography>
                        </>
                      )}
                    </Box>
                  </Grid>
                </HtmlTooltip>
              </Grid>
            </CardContent>
          </Card>
        </Fade>
      )}
    </>
  );
};
