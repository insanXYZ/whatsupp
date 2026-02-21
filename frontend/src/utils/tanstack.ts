import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { API } from "./axios";
import { ToastError, ToastSuccess } from "./toast";

export const useQueryData = (keys: any[], url: string) => {
  return useQuery({
    queryKey: keys,
    queryFn: async () => {
      const res = await API.get(url);
      return res.data;
    },
  });
};

export enum HttpMethod {
  POST = "POST",
  PUT = "PUT",
  DELETE = "DELETE",
  GET = "GET",
}

export enum ContentType {
  JSON = "application/json",
  FORM = "multipart/form-data",
}

interface Mutate {
  url: string;
  body: any;
  method: HttpMethod;
  contentType?: ContentType;
}

export function Mutation(mutationKey: any[], useToast: boolean = false) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      url,
      body,
      method,
      contentType = ContentType.JSON,
    }: Mutate) => {
      const res = await API({
        url: url,
        data: body,
        method: method,
        headers: {
          "Content-Type": contentType,
        },
      });

      return res.data;
    },
    mutationKey,
    onSuccess: (data: any) => {
      if (data.message && useToast) {
        ToastSuccess(data.message);
      }

      return queryClient.invalidateQueries({ queryKey: mutationKey });
    },
    onError: (err: any) => {
      if (err.response?.data && useToast) {
        ToastError(err.response.data.message);
      }
    },
  });
}
