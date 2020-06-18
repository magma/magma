// flow-declare typed signature: a37ab6db80ca6018833817a672baa71b
// flow-declare typed version: <<STUB>>/@turf/turf_v5.x.x

/**
 * @format
 */

declare module '@turf/turf' {
  declare export type GeoGeometryType =
    | 'Point'
    | 'MultiPoint'
    | 'LineString'
    | 'MultiLineString'
    | 'Polygon'
    | 'MultiPolygon'
    | 'GeometryCollection';

  declare export type JsonObj = {[string]: *};
  declare export type GeoCoord = [number, number] | [number, number, number];

  declare export type GeoGeometry = {|
    type: GeoGeometryType,
    coordinates: Array<GeoCoord>,
    properties: JsonObj,
  |};

  declare export type GeoFeature = {|
    type: 'Feature',
    geometry: GeoGeometry,
    properties: JsonObj,
    id?: number,
  |};

  declare export type GeoFeatureCollection = {|
    type: 'FeatureCollection',
    features: Array<GeoFeature>,
    properties: JsonObj,
  |};

  declare export type GeoJson = GeoFeature | GeoFeatureCollection | GeoGeometry;

  declare export function buffer(
    GeoFeatureCollection | GeoFeature | GeoGeometry,
    number,
    ?{units?: string},
  ): GeoFeature;
  declare export function convex(GeoFeatureCollection): GeoGeometry;
  declare export function featureCollection(
    Array<GeoFeature>,
  ): GeoFeatureCollection;
  declare export function feature(GeoGeometry): GeoFeature;
  declare export function point(
    [number, number, ?number],
    ?JsonObj,
    ?{bbox?: ?Array<GeoCoord>, id?: string | number},
  ): GeoFeature;
  declare export function transformRotate(
    GeoJson,
    number,
    ?{mutate?: boolean, pivot?: 'centroid' | GeoCoord},
  ): GeoJson;
  declare export function transformTranslate(
    GeoJson,
    number,
    number,
    ?{mutate?: boolean, units?: 'kilometers', zTranslation?: number},
  ): GeoJson;

  declare export function bearing(GeoCoord, GeoCoord): number;
}
