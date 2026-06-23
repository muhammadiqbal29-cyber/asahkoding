import { NextResponse } from 'next/server';

async function handleProxy(request, { params }) {
  const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080';
  
  // params.route is an array of path segments, e.g. ['auth', 'login']
  const routePath = params.route ? params.route.join('/') : '';
  const url = `${backendUrl}/api/${routePath}`;

  try {
    // We only pass headers and body for methods that allow them
    const fetchOptions = {
      method: request.method,
      headers: {
        'Content-Type': request.headers.get('content-type') || 'application/json',
      },
    };

    if (request.headers.has('authorization')) {
      fetchOptions.headers['Authorization'] = request.headers.get('authorization');
    }

    if (request.method !== 'GET' && request.method !== 'HEAD') {
      fetchOptions.body = await request.text();
    }

    const response = await fetch(url, fetchOptions);
    
    // Attempt to parse JSON response
    let data;
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      data = await response.json();
    } else {
      data = await response.text();
    }

    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('BFF Proxy Error:', error);
    return NextResponse.json(
      { error: 'Internal Server Error', message: error.message },
      { status: 500 }
    );
  }
}

export const GET = handleProxy;
export const POST = handleProxy;
export const PUT = handleProxy;
export const DELETE = handleProxy;
export const PATCH = handleProxy;
