/* @meta
{
  "name": "weixin/test-debug",
  "description": "Debug test",
  "domain": "mp.weixin.qq.com",
  "args": {"url": {"required": true}},
  "runtime": "http",
  "readOnly": true
}
*/

async function(args) {
  const resp = await fetch(args.url, {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
      'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
      'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8'
    }
  });
  return {status: resp.status, ok: resp.ok};
}
